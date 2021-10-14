package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestGetAllChambers(t *testing.T) {
	t.Parallel()
	t.Run("getAllChambers", getAllChambers)
	t.Run("getAllChambersEmpty", getAllChambersEmpty)
	t.Run("getAllChambersRepoError", getAllChambersRepoError)
	t.Run("getAllChambersRespondError", getAllChambersRespondError)
}

func getAllChambers(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupRequest(http.MethodGet, "/chamberss", nil)

	expected := []storage.Chamber{
		{ID: chamberID},
		{ID: "dd2610fe-95fc-45f3-8dd8-3051fb1bd4c1"},
	}

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(expected, nil)

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []storage.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func getAllChambersEmpty(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupRequest(http.MethodGet, "/chamberss", nil)

	expected := []storage.Chamber{}

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(expected, nil)

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []storage.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func getAllChambersRepoError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupRequest(http.MethodGet, "/chamberss", nil)

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return([]storage.Chamber{}, errors.New("repoMock error"))

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	// TODO: Use ErrorContains() when a release contain this PR is tagged: https://github.com/stretchr/testify/pull/1022
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "get all chambers from"))
}

func getAllChambersRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupRequest(http.MethodGet, "/chamberss", nil)

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return([]storage.Chamber{}, nil)

	handler := &handlers.ChambersHandler{Repo: repoMock}
	// use new ctx to force error
	err := handler.GetAll(context.Background(), w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}

func TestGetChamber(t *testing.T) {
	t.Parallel()
	t.Run("getChamberFound", getChamberFound)
	t.Run("getChamberNotFound", getChamberNotFound)
	t.Run("getChamberRepoError", getChamberRepoError)
	t.Run("getChamberRespondError", getChamberRespondError)
}

func getChamberFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupRequest(http.MethodGet, "/chambers/"+chamberID, nil)

	expected := storage.Chamber{ID: chamberID}
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Get", chamberID).Return(&expected, nil)

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	chamber := storage.Chamber{}
	err = json.Unmarshal(bodyBytes, &chamber)
	assert.NoError(t, err)
	assert.Equal(t, expected, chamber)
}

func getChamberNotFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupRequest(http.MethodGet, "/chambers/"+chamberID, nil)

	var expected *storage.Chamber

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Get", chamberID).Return(expected, nil)

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))
}

func getChamberRepoError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupRequest(http.MethodGet, "/chambers/"+chamberID, nil)

	var expected *storage.Chamber

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Get", chamberID).Return(expected, errors.New("repoMock error"))

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "get chamber from"))
}

func getChamberRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupRequest(http.MethodGet, "/chambers/"+chamberID, nil)

	expected := storage.Chamber{ID: chamberID}
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Get", chamberID).Return(&expected, nil)

	handler := &handlers.ChambersHandler{Repo: repoMock}
	// use new ctx to force error
	err := handler.Get(context.Background(), w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}

func TestSaveChamber(t *testing.T) {
	t.Parallel()
	t.Run("saveChamber", saveChamber)
	t.Run("saveChamberParseError", saveChamberParseError)
	t.Run("saveChamberRepoError", saveChamberRepoError)
	t.Run("saveChamberRespondError", saveChamberRespondError)
}

func saveChamber(t *testing.T) {
	t.Parallel()

	chamber := storage.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(chamber)
	w, r, ctx := setupRequest(http.MethodPost, "/chambers", bytes.NewBuffer(jsonBytes))

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Save", &chamber).Return(nil)

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := storage.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
}

func saveChamberParseError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupRequest(http.MethodPost, "/chambers", nil)

	repoMock := &mocks.ChamberRepo{}

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), parseErrorMsg)
}

func saveChamberRepoError(t *testing.T) {
	t.Parallel()

	chamber := storage.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(chamber)
	w, r, ctx := setupRequest(http.MethodPost, "/chambers", bytes.NewBuffer(jsonBytes))

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Save", &chamber).Return(errors.New("repoMock error"))

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "save chamber to"))
}

func saveChamberRespondError(t *testing.T) {
	t.Parallel()

	chamber := storage.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(chamber)
	w, r, _ := setupRequest(http.MethodPost, "/chambers", bytes.NewBuffer(jsonBytes))

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Save", &chamber).Return(nil)

	handler := &handlers.ChambersHandler{Repo: repoMock}
	// use new ctx to force error
	err := handler.Save(context.Background(), w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}

func TestDeleteChamber(t *testing.T) {
	t.Parallel()
	t.Run("deleteChamberFound", deleteChamberFound)
	t.Run("deleteChamberNotFound", deleteChamberNotFound)
	t.Run("deleteChamberRepoGetError", deleteChamberRepoGetError)
	t.Run("deleteChamberRepoDeleteError", deleteChamberRepoDeleteError)
	t.Run("deleteChamberRespondError", deleteChamberRespondError)
}

func deleteChamberFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupRequest(http.MethodDelete, "/chambers/"+chamberID, nil)

	chamber := storage.Chamber{ID: chamberID}
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Get", chamberID).Return(&chamber, nil)
	repoMock.On("Delete", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
}

func deleteChamberNotFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupRequest(http.MethodDelete, "/chambers/"+chamberID, nil)

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Get", chamberID).Return(nil, nil)

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))
}

func deleteChamberRepoGetError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupRequest(http.MethodDelete, "/chambers/"+chamberID, nil)

	var chamber *storage.Chamber

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Get", chamberID).Return(chamber, errors.New("repoMock error"))

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "get chamber from"))
}

func deleteChamberRepoDeleteError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupRequest(http.MethodDelete, "/chambers/"+chamberID, nil)

	chamber := storage.Chamber{ID: chamberID}
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Get", chamberID).Return(&chamber, nil)
	repoMock.On("Delete", chamberID).Return(errors.New("repoMock error"))

	handler := &handlers.ChambersHandler{Repo: repoMock}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "delete chamber from"))
}

func deleteChamberRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupRequest(http.MethodDelete, "/chambers/"+chamberID, nil)

	chamber := storage.Chamber{ID: chamberID}
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("Get", chamberID).Return(&chamber, nil)
	repoMock.On("Delete", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{Repo: repoMock}
	// use new ctx to force error
	err := handler.Delete(context.Background(), w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}
