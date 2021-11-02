package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/controller"
	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

	w, r, ctx := setupHandlerTest("", nil)
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

	w, r, ctx := setupHandlerTest("", nil)
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

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetAllChambers").Return([]chamber.Chamber{}, errSomeError)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	// TODO: Waiting on PR for ErrorContains(): https://github.com/stretchr/testify/pull/1022
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "get all chambers from"))
}

func getAllChambersRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest("", nil)
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

	w, r, ctx := setupHandlerTest("", nil)
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

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	var expected *chamber.Chamber

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(expected, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func getChamberControllerError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	var expected *chamber.Chamber

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(expected, errSomeError)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "get chamber "+chamberID+" from"))
}

func getChamberRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest("", nil)
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
	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
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

	w, r, ctx := setupHandlerTest("", nil)
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
	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
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
	w, r, _ := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
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
	t.Run("deleteChamberFermenting", deleteChamberFermenting)
	t.Run("deleteChamberNotFound", deleteChamberNotFound)
	t.Run("deleteChamberControllerGetError", deleteChamberControllerGetError)
	t.Run("deleteChamberControllerStopFermentationNotFermentingError",
		deleteChamberControllerStopFermentationNotFermentingError)
	t.Run("deleteChamberControllerStopFermentationError", deleteChamberControllerStopFermentationError)
	t.Run("deleteChamberControllerDeleteError", deleteChamberControllerDeleteError)
	t.Run("deleteChamberRespondError", deleteChamberRespondError)
}

func deleteChamberFermenting(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, hook := logtest.NewNullLogger()

	chamber := chamber.Chamber{ID: chamberID}
	_ = chamber.Configure(ctx, l)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(&chamber, nil)
	controllerMock.On("StopFermentation", chamberID).Return(nil)
	controllerMock.On("DeleteChamber", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
	assert.Nil(t, hook.LastEntry())
}

func deleteChamberNotFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(nil, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func deleteChamberControllerGetError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	var chamber *chamber.Chamber

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(chamber, errSomeError)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "get chamber "+chamberID+" from"))
}

func deleteChamberControllerStopFermentationNotFermentingError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, hook := logtest.NewNullLogger()

	c := chamber.Chamber{ID: chamberID}
	_ = c.Configure(ctx, l)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(&c, nil)
	controllerMock.On("StopFermentation", chamberID).Return(errors.Wrap(chamber.ErrNotFermenting, "some error"))
	controllerMock.On("DeleteChamber", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
	assert.Equal(t, hook.LastEntry().Level, logrus.WarnLevel)
	assert.Equal(t, hook.LastEntry().Message, "Error occurred while stopping fermentation")
}

func deleteChamberControllerStopFermentationError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	chamber := chamber.Chamber{ID: chamberID}
	_ = chamber.Configure(ctx, l)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(&chamber, nil)
	controllerMock.On("StopFermentation", chamberID).Return(errSomeError)
	controllerMock.On("DeleteChamber", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf("could not stop fermentation for chamber %s", chamberID))
}

func deleteChamberControllerDeleteError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	chamber := chamber.Chamber{ID: chamberID}
	_ = chamber.Configure(ctx, l)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(&chamber, nil)
	controllerMock.On("DeleteChamber", chamberID).Return(errSomeError)
	controllerMock.On("StopFermentation", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "delete chamber "+chamberID+" from"))
}

func deleteChamberRespondError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	chamber := chamber.Chamber{ID: chamberID}
	_ = chamber.Configure(ctx, l)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("GetChamber", chamberID).Return(&chamber, nil)
	controllerMock.On("DeleteChamber", chamberID).Return(nil)
	controllerMock.On("StopFermentation", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	// use new ctx to force error
	err := handler.Delete(context.Background(), w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestStartFermentation(t *testing.T) {
	t.Parallel()
	t.Run("startFermentationSuccess", startFermentationSuccess)
	t.Run("startFermentationStepParseError", startFermentationStepParseError)
	t.Run("startFermentationInvalidStep", startFermentationInvalidStep)
	t.Run("startFermentationNotFound", startFermentationNotFound)
	t.Run("startFermentationNoBatch", startFermentationNoBatch)
	t.Run("startFermentationChamberControllerError", startFermentationChamberControllerError)
	t.Run("startFermentationRespondError", startFermentationRespondError)
}

func startFermentationSuccess(t *testing.T) {
	t.Parallel()

	step := 1

	w, r, ctx := setupHandlerTest("step="+strconv.Itoa(step), nil)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("StartFermentation", chamberID, step).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
}

func startFermentationStepParseError(t *testing.T) {
	t.Parallel()

	step := "One"
	w, r, ctx := setupHandlerTest("step="+step, nil)
	controllerMock := &mocks.ChamberController{}

	handler := &handlers.ChambersHandler{ChamberController: controllerMock}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf("step %s is invalid", step))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func startFermentationInvalidStep(t *testing.T) {
	t.Parallel()

	step := 2
	w, r, ctx := setupHandlerTest("step="+strconv.Itoa(step), nil)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("StartFermentation", chamberID, step).Return(errors.Wrap(chamber.ErrInvalidStep, "some error"))

	handler := &handlers.ChambersHandler{ChamberController: controllerMock}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf("step %d is invalid for chamber '%s'", step, chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func startFermentationNotFound(t *testing.T) {
	t.Parallel()

	step := 1
	w, r, ctx := setupHandlerTest("step="+strconv.Itoa(step), nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("StartFermentation", chamberID, step).Return(controller.ErrNotFound)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func startFermentationNoBatch(t *testing.T) {
	t.Parallel()

	step := 1
	w, r, ctx := setupHandlerTest("step="+strconv.Itoa(step), nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("StartFermentation", chamberID, step).Return(errors.Wrap(chamber.ErrNoCurrentBatch, "some error"))

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf("chamber '%s' does not have a current batch", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func startFermentationChamberControllerError(t *testing.T) {
	t.Parallel()

	step := 1
	w, r, ctx := setupHandlerTest("step="+strconv.Itoa(step), nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("StartFermentation", chamberID, step).Return(errSomeError)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf("could not start fermentation for chamber %s", chamberID))
}

func startFermentationRespondError(t *testing.T) {
	t.Parallel()

	step := 1
	w, r, ctx := setupHandlerTest("step="+strconv.Itoa(step), nil)
	l, _ := logtest.NewNullLogger()

	chamber := chamber.Chamber{ID: chamberID}
	_ = chamber.Configure(ctx, l)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("StartFermentation", chamberID, step).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	// use new ctx to force error
	err := handler.Start(context.Background(), w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestStopFermentation(t *testing.T) {
	t.Parallel()
	t.Run("stopFermentationSuccess", stopFermentationSuccess)
	t.Run("stopFermentationNotFound", stopFermentationNotFound)
	t.Run("stopFermentationNotFermenting", stopFermentationNotFermenting)
	t.Run("stopFermentationChamberControllerError", stopFermentationChamberControllerError)
	t.Run("stopFermentationRespondError", stopFermentationRespondError)
}

func stopFermentationSuccess(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("StopFermentation", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock}
	err := handler.Stop(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
}

func stopFermentationNotFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("StopFermentation", chamberID).Return(controller.ErrNotFound)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Stop(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func stopFermentationNotFermenting(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("StopFermentation", chamberID).Return(errors.Wrap(chamber.ErrNotFermenting, "some error"))

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Stop(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf("chamber '%s' is not fermenting", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func stopFermentationChamberControllerError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.ChamberController{}
	controllerMock.On("StopFermentation", chamberID).Return(errSomeError)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Stop(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf("could not stop fermentation for chamber %s", chamberID))
}

func stopFermentationRespondError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	chamber := chamber.Chamber{ID: chamberID}
	_ = chamber.Configure(ctx, l)
	controllerMock := &mocks.ChamberController{}
	controllerMock.On("StopFermentation", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	// use new ctx to force error
	err := handler.Stop(context.Background(), w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}
