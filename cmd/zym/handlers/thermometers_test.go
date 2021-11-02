package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetAllThermometers(t *testing.T) {
	t.Parallel()
	t.Run("getAllThermometers", getAllThermometers)
	t.Run("getAllThermometersEmpty", getAllThermometersEmpty)
	t.Run("getAllThermometersRepoError", getAllThermometersRepoError)
	t.Run("getAllThermometersRespondError", getAllThermometersRespondError)
}

func getAllThermometers(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	expected := []string{"28-0000071cbc72", "28-0000041ab222"}

	repoMock := &mocks.ThermometerRepo{}
	repoMock.On("GetThermometerIDs").Return(expected, nil)

	handler := &handlers.ThermometersHandler{Repo: repoMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func getAllThermometersEmpty(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	expected := []string{}
	repoMock := &mocks.ThermometerRepo{}
	repoMock.On("GetThermometerIDs").Return(expected, nil)

	handler := &handlers.ThermometersHandler{Repo: repoMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func getAllThermometersRepoError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	repoMock := &mocks.ThermometerRepo{}
	repoMock.On("GetThermometerIDs").Return([]string{}, errSomeError)

	handler := &handlers.ThermometersHandler{Repo: repoMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	// TODO: Waiting on PR for ErrorContains(): https://github.com/stretchr/testify/pull/1022
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "get all thermometers from"))
}

func getAllThermometersRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest("", nil)
	ctx := context.Background()

	repoMock := &mocks.ThermometerRepo{}
	repoMock.On("GetThermometerIDs").Return([]string{}, nil)

	handler := &handlers.ThermometersHandler{Repo: repoMock}
	// use new ctx to force error
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}
