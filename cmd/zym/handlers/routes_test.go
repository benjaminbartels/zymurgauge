package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/benjaminbartels/zymurgauge/internal/test/stubs"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const devicePath = "/"

func TestRoutes(t *testing.T) {
	t.Parallel()

	type test struct {
		path   string
		method string
		body   interface{}
		code   int
	}

	testCases := []test{
		{path: "/v1/chambers", method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/v1/chambers/" + chamberID, method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/v1/chambers", method: http.MethodPost, body: &chamber.Chamber{ID: chamberID}, code: http.StatusOK},
		{path: "/v1/chambers/" + chamberID, method: http.MethodDelete, body: nil, code: http.StatusOK},
		{path: "/v1/chambers/" + chamberID + "/start?step=Primary", method: http.MethodPost, body: nil, code: http.StatusOK},
		{path: "/v1/chambers/" + chamberID + "/stop", method: http.MethodPost, body: nil, code: http.StatusOK},
		{path: "/v1/thermometers", method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/v1/batches", method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/v1/batches/" + batchID, method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/v1/bad_path/" + batchID, method: http.MethodGet, body: nil, code: http.StatusNotFound},
	}

	for _, tc := range testCases {
		tc := tc

		ctx := context.Background()
		l, _ := logtest.NewNullLogger()

		c := &chamber.Chamber{
			ID: chamberID,
			DeviceConfig: chamber.DeviceConfig{
				ChillerGPIO:         "2",
				HeaterGPIO:          "3",
				BeerThermometerType: "ds18b20",
				BeerThermometerID:   "1",
			},
			CurrentBatch: &brewfather.BatchDetail{
				Fermentation: brewfather.Fermentation{
					Steps: []brewfather.FermentationStep{
						{
							Type:            "Primary",
							StepTemperature: 22,
						},
						{
							Type:            "Secondary",
							StepTemperature: 20,
						},
					},
				},
			},
		}

		configuratorMock := &mocks.Configurator{}
		configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
		configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
		configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

		r := brewfather.BatchDetail{ID: batchID}

		serviceMock := &mocks.Service{}
		serviceMock.On("GetAllSummaries", mock.Anything).Return([]brewfather.BatchSummary{}, nil)
		serviceMock.On("GetDetail", mock.Anything, batchID).Return(&r, nil)

		err := c.Configure(configuratorMock, serviceMock, false, l)
		assert.NoError(t, err)

		chambers := []*chamber.Chamber{c}

		controllerMock := &mocks.Controller{}
		controllerMock.On("GetAll").Return(chambers, nil)
		controllerMock.On("Get", mock.Anything).Return(c, nil)
		controllerMock.On("Save", mock.Anything).Return(nil)
		controllerMock.On("Delete", mock.Anything).Return(nil)
		controllerMock.On("StartFermentation", chamberID, "Primary").Return(nil)
		controllerMock.On("StopFermentation", chamberID).Return(nil)

		shutdown := make(chan os.Signal, 1)
		logger, _ := logtest.NewNullLogger()

		app := handlers.NewAPI(controllerMock, devicePath, serviceMock, shutdown, logger)

		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()

			if tc.path == "/v1/chambers/"+chamberID+"/stop" {
				c, err := controllerMock.Get(chamberID)
				assert.NoError(t, err)

				err = c.StartFermentation(ctx, "Primary")
				assert.NoError(t, err)
			}

			w := httptest.NewRecorder()
			jsonBytes, _ := json.Marshal(tc.body)
			r := httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(jsonBytes))
			app.ServeHTTP(w, r)
			assert.Equal(t, tc.code, w.Code)
		})
	}
}

func TestDebugMux(t *testing.T) {
	t.Parallel()

	mux := handlers.DebugMux()

	type test struct {
		path   string
		method string
		code   int
	}

	testCases := []test{
		{path: "/debug/pprof/", method: http.MethodGet, code: http.StatusOK},
		{path: "/debug/pprof/cmdline", method: http.MethodGet, code: http.StatusOK},
		{path: "/debug/pprof/profile?seconds=1", method: http.MethodGet, code: http.StatusOK},
		{path: "/debug/pprof/symbol", method: http.MethodGet, code: http.StatusOK},
		{path: "/debug/pprof/trace", method: http.MethodGet, code: http.StatusOK},
		{path: "/debug/vars", method: http.MethodGet, code: http.StatusOK},
		{path: "/debug/pprof/bad_path", method: http.MethodGet, code: http.StatusNotFound},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.path, nil)
			mux.ServeHTTP(w, r)
			assert.Equal(t, tc.code, w.Code)
		})
	}
}
