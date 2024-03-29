package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/auth"
	"github.com/benjaminbartels/zymurgauge/internal/batch"
	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/platform/debug"
	"github.com/benjaminbartels/zymurgauge/internal/settings"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/benjaminbartels/zymurgauge/internal/test/stubs"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	devicePath            = "/"
	readingUpdateInterval = 100 * time.Millisecond
)

func TestRoutes(t *testing.T) {
	t.Parallel()

	type test struct {
		path   string
		method string
		body   interface{}
		code   int
	}

	testCases := []test{
		{path: "/api/v1/chambers", method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/api/v1/chambers/" + chamberID, method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/api/v1/chambers", method: http.MethodPost, body: &chamber.Chamber{ID: chamberID}, code: http.StatusOK},
		{path: "/api/v1/chambers/" + chamberID, method: http.MethodDelete, body: nil, code: http.StatusOK},
		{path: "/api/v1/chambers/" + chamberID + "/start?step=A", method: http.MethodPost, body: nil, code: http.StatusOK},
		{path: "/api/v1/chambers/" + chamberID + "/stop", method: http.MethodPost, body: nil, code: http.StatusOK},
		{path: "/api/v1/thermometers", method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/api/v1/batches", method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/api/v1/batches/" + batchID, method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/api/v1/bad_path/" + batchID, method: http.MethodGet, body: nil, code: http.StatusNotFound},
		{path: "/index.html", method: http.MethodGet, body: nil, code: http.StatusOK},
	}

	for _, tc := range testCases {
		tc := tc

		ctx := context.Background()
		l, _ := logtest.NewNullLogger()
		m := &mocks.Metrics{}
		m.On("Gauge", mock.Anything, mock.Anything).Return()

		c := &chamber.Chamber{
			ID: chamberID,
			DeviceConfig: chamber.DeviceConfig{
				ChillerGPIO:         "2",
				HeaterGPIO:          "3",
				BeerThermometerType: "ds18b20",
				BeerThermometerID:   "1",
			},
			CurrentBatch: &batch.Detail{
				Recipe: batch.Recipe{
					Fermentation: batch.Fermentation{
						Steps: []batch.FermentationStep{
							{
								Name:        "A",
								Temperature: 22,
							},
							{
								Name:        "B",
								Temperature: 20,
							},
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
		serviceMock.On("GetAllBatchSummaries", mock.Anything).Return([]brewfather.BatchSummary{}, nil)
		serviceMock.On("GetBatchDetail", mock.Anything, batchID).Return(&r, nil)
		serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil)

		err := c.Configure(configuratorMock, serviceMock, l, m, readingUpdateInterval)
		assert.NoError(t, err)

		chambers := []*chamber.Chamber{c}

		controllerMock := &mocks.Controller{}
		controllerMock.On("GetAll").Return(chambers, nil)
		controllerMock.On("Get", mock.Anything).Return(c, nil)
		controllerMock.On("Save", mock.Anything).Return(nil)
		controllerMock.On("Delete", mock.Anything).Return(nil)
		controllerMock.On("StartFermentation", chamberID, "A").Return(nil)
		controllerMock.On("StopFermentation", chamberID).Return(nil)

		s := &settings.Settings{
			AppSettings: settings.AppSettings{AuthSecret: "my-auth-secret"},
		}

		settingsMock := &mocks.SettingsRepo{}
		settingsMock.On("Get").Return(s, nil)

		shutdown := make(chan os.Signal, 1)
		logger, _ := logtest.NewNullLogger()

		fsMock := &mocks.FileReader{}
		fsMock.On("ReadFile", "build/index.html").Return([]byte(""), nil)

		app, _ := handlers.NewApp(controllerMock, devicePath, serviceMock, settingsMock, nil, fsMock,
			shutdown, logger)

		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()

			if tc.path == "/v1/chambers/"+chamberID+"/stop" {
				c, err := controllerMock.Get(chamberID)
				assert.NoError(t, err)

				err = c.StartFermentation(ctx, "A")
				assert.NoError(t, err)
			}

			w := httptest.NewRecorder()
			jsonBytes, _ := json.Marshal(tc.body)
			r := httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(jsonBytes))
			token, _ := auth.CreateToken("my-auth-secret", "username", 1*time.Minute)

			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			app.ServeHTTP(w, r)
			assert.Equal(t, tc.code, w.Code)
		})
	}
}

func TestDebugMux(t *testing.T) {
	t.Parallel()

	mux := debug.Mux()

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
