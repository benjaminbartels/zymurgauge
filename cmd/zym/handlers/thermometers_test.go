package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest // False positives with r.Run not in a loop
func TestGetAllThermometers(t *testing.T) {
	t.Parallel()
	t.Run("getAllThermometers", getAllThermometers)
	t.Run("getAllThermometersEmpty", getAllThermometersEmpty)
	t.Run("getAllThermometersRepoError", getAllThermometersRepoError)
	t.Run("getAllThermometersRespondError", getAllThermometersRespondError)
}

func getAllThermometers(t *testing.T) {
	t.Parallel()

	file, err := os.Create("28-1234567890")
	assert.NoError(t, err)
	file.Close()

	defer os.Remove(file.Name())

	w, r, ctx := setupHandlerTest("", nil)

	expected := []string{filepath.Base(file.Name())}

	handler := &handlers.ThermometersHandler{DevicePath: "./"}
	err = handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := io.ReadAll(resp.Body)
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

	handler := &handlers.ThermometersHandler{DevicePath: devicePath}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	result := []string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func getAllThermometersRepoError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)

	handler := &handlers.ThermometersHandler{DevicePath: "]["} // force glob to error
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	// TODO: Waiting on PR for ErrorContains(): https://github.com/stretchr/testify/pull/1022
	assert.Contains(t, err.Error(), "could not get all thermometers ids from onewire bus")
}

func getAllThermometersRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest("", nil)
	ctx := context.Background()

	handler := &handlers.ThermometersHandler{DevicePath: devicePath}
	// use new ctx to force error
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}
