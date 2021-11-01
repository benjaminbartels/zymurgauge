package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/julienschmidt/httprouter"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetAllChambers(t *testing.T) {
	t.Parallel()
	t.Run("getAllChambers", getAllChambers)
	t.Run("getAllChambersEmpty", getAllChambersEmpty)
	t.Run("getAllChambersRepoError", getAllChambersRepoError)
	t.Run("getAllChambersRespondError", getAllChambersRespondError)
}

func getAllChambers(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()
	expected := []chamber.Chamber{
		{ID: chamberID},
		{ID: "dd2610fe-95fc-45f3-8dd8-3051fb1bd4c1"},
	}

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetAllChambers").Return(expected, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []chamber.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func getAllChambersEmpty(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()
	expected := []chamber.Chamber{}

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetAllChambers").Return(expected, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []chamber.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func getAllChambersRepoError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetAllChambers").Return([]chamber.Chamber{}, errDeadDatabase)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	// TODO: Waiting on PR for ErrorContains(): https://github.com/stretchr/testify/pull/1022
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "get all chambers from"))
}

func getAllChambersRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetAllChambers").Return([]chamber.Chamber{}, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	// use new ctx to force error
	err := handler.GetAll(context.Background(), w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetChamber(t *testing.T) {
	t.Parallel()
	t.Run("getChamberFound", getChamberFound)
	t.Run("getChamberNotFound", getChamberNotFound)
	t.Run("getChamberControllerError", getChamberControllerError)
	t.Run("getChamberRespondError", getChamberRespondError)
}

func getChamberFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()
	expected := chamber.Chamber{ID: chamberID}
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(&expected, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	chamber := chamber.Chamber{}
	err = json.Unmarshal(bodyBytes, &chamber)
	assert.NoError(t, err)
	assert.Equal(t, expected, chamber)
}

func getChamberNotFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()

	var expected *chamber.Chamber

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(expected, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))
}

func getChamberControllerError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()

	var expected *chamber.Chamber

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(expected, errDeadDatabase)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "get chamber "+chamberID+" from"))
}

func getChamberRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()

	expected := chamber.Chamber{ID: chamberID}
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(&expected, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	// use new ctx to force error
	err := handler.Get(context.Background(), w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestSaveChamber(t *testing.T) {
	t.Parallel()
	t.Run("saveChamber", saveChamber)
	t.Run("saveChamberParseError", saveChamberParseError)
	t.Run("saveChamberControllerError", saveChamberControllerError)
	t.Run("saveChamberRespondError", saveChamberRespondError)
}

func saveChamber(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(c)
	w, r, ctx := setupHandlerTest(bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("SaveChamber", &c).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := chamber.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
}

func saveChamberParseError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), parseErrorMsg)
}

func saveChamberControllerError(t *testing.T) {
	t.Parallel()

	chamber := chamber.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(chamber)
	w, r, ctx := setupHandlerTest(bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("SaveChamber", &chamber).Return(errors.New("controllerMock error"))

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "save chamber to"))
}

func saveChamberRespondError(t *testing.T) {
	t.Parallel()

	chamber := chamber.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(chamber)
	w, r, _ := setupHandlerTest(bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("SaveChamber", &chamber).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	// use new ctx to force error
	err := handler.Save(context.Background(), w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestDeleteChamber(t *testing.T) {
	t.Parallel()
	t.Run("deleteChamberFound", deleteChamberFound)
	t.Run("deleteChamberNotFound", deleteChamberNotFound)
	t.Run("deleteChamberControllerGetError", deleteChamberControllerGetError)
	t.Run("deleteChamberControllerDeleteError", deleteChamberControllerDeleteError)
	t.Run("deleteChamberRespondError", deleteChamberRespondError)
}

func deleteChamberFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()

	chamber := chamber.Chamber{ID: chamberID}
	_ = chamber.Configure(ctx, l)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(&chamber, nil)
	controllerMock.On("DeleteChamber", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
}

func deleteChamberNotFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(nil, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))
}

func deleteChamberControllerGetError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()

	var chamber *chamber.Chamber

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(chamber, errDeadDatabase)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "get chamber "+chamberID+" from"))
}

func deleteChamberControllerDeleteError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()

	chamber := chamber.Chamber{ID: chamberID}
	_ = chamber.Configure(ctx, l)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(&chamber, nil)
	controllerMock.On("DeleteChamber", chamberID).Return(errDeadDatabase)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "delete chamber "+chamberID+" from"))
}

func deleteChamberRespondError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)
	l, _ := logtest.NewNullLogger()

	chamber := chamber.Chamber{ID: chamberID}
	_ = chamber.Configure(ctx, l)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(&chamber, nil)
	controllerMock.On("DeleteChamber", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	// use new ctx to force error
	err := handler.Delete(context.Background(), w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}
