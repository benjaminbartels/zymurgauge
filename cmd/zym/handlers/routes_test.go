package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoutes(t *testing.T) {
	t.Parallel()

	chamber := storage.Chamber{ID: chamberID}
	recipe := brewfather.Recipe{ID: recipeID}

	chamberRepoMock := &mocks.ChamberRepo{}
	chamberRepoMock.On("GetAll").Return([]storage.Chamber{}, nil)
	chamberRepoMock.On("Get", mock.Anything).Return(&chamber, nil)
	chamberRepoMock.On("Save", mock.Anything).Return(nil)
	chamberRepoMock.On("Delete", mock.Anything).Return(nil)

	recipeMock := &mocks.RecipeRepo{}
	recipeMock.On("GetRecipes", mock.Anything).Return([]brewfather.Recipe{}, nil)
	recipeMock.On("GetRecipe", mock.Anything, recipeID).Return(&recipe, nil)

	shutdown := make(chan os.Signal, 1)
	logger, _ := logtest.NewNullLogger()

	app := handlers.NewAPI(chamberRepoMock, recipeMock, shutdown, logger)

	type test struct {
		name   string
		method string
		path   string
		body   interface{}
	}

	testCases := []test{
		{name: "GetAllChambers", method: http.MethodGet, path: "/chambers", body: nil},
		{name: "GetChamber", method: http.MethodGet, path: "/chambers/" + chamberID, body: nil},
		{name: "SaveChamber", method: http.MethodPost, path: "/chambers", body: chamber},
		{name: "DeleteChamber", method: http.MethodDelete, path: "/chambers/" + chamberID, body: nil},
		{name: "GetAllRecipes", method: http.MethodGet, path: "/recipes", body: nil},
		{name: "GetRecipe", method: http.MethodGet, path: "/recipes/" + recipeID, body: nil},
	}

	for _, tc := range testCases {
		tc := tc // TODO: Remove this with new linter config
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			jsonBytes, _ := json.Marshal(tc.body)
			r := httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(jsonBytes))
			app.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
