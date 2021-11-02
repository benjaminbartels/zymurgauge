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
	"github.com/benjaminbartels/zymurgauge/internal/batch"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoutes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	l, _ := logtest.NewNullLogger()

	c := chamber.Chamber{ID: chamberID}
	_ = c.Configure(ctx, l)
	r := batch.Batch{ID: batchID}

	chamberControllerMock := &mocks.ChamberController{}
	chamberControllerMock.On("GetAllChambers").Return([]chamber.Chamber{}, nil)
	chamberControllerMock.On("GetChamber", mock.Anything).Return(&c, nil)
	chamberControllerMock.On("SaveChamber", mock.Anything).Return(nil)
	chamberControllerMock.On("DeleteChamber", mock.Anything).Return(nil)
	chamberControllerMock.On("StartFermentation", chamberID, 1).Return(nil)
	chamberControllerMock.On("StopFermentation", chamberID).Return(nil)

	thermometerMock := &mocks.ThermometerRepo{}
	thermometerMock.On("GetThermometerIDs", mock.Anything).Return([]string{}, nil)

	recipeMock := &mocks.BatchRepo{}
	recipeMock.On("GetAllBatches", mock.Anything).Return([]batch.Batch{}, nil)
	recipeMock.On("GetBatch", mock.Anything, batchID).Return(&r, nil)

	shutdown := make(chan os.Signal, 1)
	logger, _ := logtest.NewNullLogger()

	app := handlers.NewAPI(chamberControllerMock, thermometerMock, recipeMock, shutdown, logger)

	type test struct {
		path   string
		method string
		body   interface{}
		code   int
	}

	testCases := []test{
		{path: "/v1/chambers", method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/v1/chambers/" + chamberID, method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/v1/chambers", method: http.MethodPost, body: &c, code: http.StatusOK},
		{path: "/v1/chambers/" + chamberID, method: http.MethodDelete, body: nil, code: http.StatusOK},
		{path: "/v1/chambers/" + chamberID + "/start?step=1", method: http.MethodPost, body: nil, code: http.StatusOK},
		{path: "/v1/chambers/" + chamberID + "/stop", method: http.MethodPost, body: nil, code: http.StatusOK},
		{path: "/v1/thermometers", method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/v1/batches", method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/v1/batches/" + batchID, method: http.MethodGet, body: nil, code: http.StatusOK},
		{path: "/v1/bad_path/" + batchID, method: http.MethodGet, body: nil, code: http.StatusNotFound},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
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
