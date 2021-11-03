package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/controller"
	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/batch"
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
	t.Run("getAllChambersRespondError", getAllChambersRespondError)
}

func getAllChambers(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()
	chambers := []chamber.Chamber{
		{ID: chamberID},
		{ID: "dd2610fe-95fc-45f3-8dd8-3051fb1bd4c1"},
	}

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return(chambers, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)

	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []chamber.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)

	assetChamberListsAreEqual(t, chambers, result)
}

func getAllChambersEmpty(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()
	chambers := []chamber.Chamber{}

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return(chambers, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []chamber.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, chambers, result)
}

func getAllChambersRespondError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	// use new ctx to force error
	err := handler.GetAll(context.Background(), w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetChamber(t *testing.T) {
	t.Parallel()
	t.Run("getChamberFound", getChamberFound)
	t.Run("getChamberNotFound", getChamberNotFound)
	t.Run("getChamberRespondError", getChamberRespondError)
}

func getChamberFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	chambers := []chamber.Chamber{{ID: chamberID}}
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return(chambers, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}

	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)

	resp := w.Result()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := chamber.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assertChambersAreEqual(t, chambers[0], result)
}

func getChamberNotFound(t *testing.T) {
	t.Parallel()

	chambers := []chamber.Chamber{}

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return(chambers, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}

	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func getChamberRespondError(t *testing.T) {
	t.Parallel()

	chambers := []chamber.Chamber{{ID: chamberID}}

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return(chambers, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}

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

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)
	repoMock.On("SaveChamber", &c).Return(nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}

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

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}

	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), parseErrorMsg)
}

func saveChamberControllerError(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(c)

	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)
	repoMock.On("SaveChamber", &c).Return(errors.New("repoMock error"))

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}

	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "save chamber to"))
}

func saveChamberRespondError(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(c)

	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)
	repoMock.On("SaveChamber", &c).Return(nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}

	// use new ctx to force error
	err := handler.Save(context.Background(), w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestDeleteChamber(t *testing.T) {
	t.Parallel()
	t.Run("deleteChamberFermenting", deleteChamberFermenting)
	t.Run("deleteChamberNotFound", deleteChamberNotFound)
	t.Run("deleteChamberControllerDeleteError", deleteChamberControllerDeleteError)
	t.Run("deleteChamberControllerStopFermentationError", deleteChamberControllerStopFermentationError)
	t.Run("deleteChamberControllerDeleteError", deleteChamberControllerDeleteError)
	t.Run("deleteChamberRespondError", deleteChamberRespondError)
}

func deleteChamberFermenting(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{
		ID: chamberID,
		CurrentBatch: &batch.Batch{
			Fermentation: batch.Fermentation{
				Steps: []batch.FermentationStep{{StepTemp: 22}},
			},
		},
	}

	w, r, ctx := setupHandlerTest("", nil)
	l, hook := logtest.NewNullLogger()

	_ = c.Configure(ctx, l)
	_ = c.StartFermentation(1)
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)
	repoMock.On("DeleteChamber", chamberID).Return(nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
	assert.Nil(t, hook.LastEntry())
}

func deleteChamberNotFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func deleteChamberControllerDeleteError(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{ID: chamberID}
	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)
	repoMock.On("DeleteChamber", chamberID).Return(errSomeError)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "delete chamber "+chamberID+" from"))
}

func deleteChamberControllerStopFermentationError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, hook := logtest.NewNullLogger()

	c := chamber.Chamber{ID: chamberID}
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)
	repoMock.On("DeleteChamber", chamberID).Return(nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
	assert.Equal(t, hook.LastEntry().Level, logrus.WarnLevel)
	assert.Equal(t, hook.LastEntry().Message, "Error occurred while stopping fermentation")
}

func deleteChamberRespondError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	c := chamber.Chamber{ID: chamberID}
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)
	repoMock.On("DeleteChamber", chamberID).Return(nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
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
	t.Run("startFermentationRespondError", startFermentationRespondError)
}

func startFermentationSuccess(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{
		ID: chamberID,
		CurrentBatch: &batch.Batch{
			Fermentation: batch.Fermentation{
				Steps: []batch.FermentationStep{{StepTemp: 22}},
			},
		},
	}
	step := 1
	w, r, ctx := setupHandlerTest("step="+strconv.Itoa(step), nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
}

func startFermentationStepParseError(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{
		ID: chamberID,
		CurrentBatch: &batch.Batch{
			Fermentation: batch.Fermentation{
				Steps: []batch.FermentationStep{{StepTemp: 22}},
			},
		},
	}
	step := "One"
	w, r, ctx := setupHandlerTest("step="+step, nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf("step %s is invalid", step))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func startFermentationInvalidStep(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{
		ID: chamberID,
		CurrentBatch: &batch.Batch{
			Fermentation: batch.Fermentation{
				Steps: []batch.FermentationStep{{StepTemp: 22}},
			},
		},
	}
	step := 2
	w, r, ctx := setupHandlerTest("step="+strconv.Itoa(step), nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager}
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

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func startFermentationNoBatch(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{ID: chamberID}
	step := 1
	w, r, ctx := setupHandlerTest("step="+strconv.Itoa(step), nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf("chamber '%s' does not have a current batch", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func startFermentationRespondError(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{
		ID: chamberID,
		CurrentBatch: &batch.Batch{
			Fermentation: batch.Fermentation{
				Steps: []batch.FermentationStep{{StepTemp: 22}},
			},
		},
	}
	step := 1
	w, r, ctx := setupHandlerTest("step="+strconv.Itoa(step), nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
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
	t.Run("stopFermentationRespondError", stopFermentationRespondError)
}

func stopFermentationSuccess(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{
		ID: chamberID,
		CurrentBatch: &batch.Batch{
			Fermentation: batch.Fermentation{
				Steps: []batch.FermentationStep{{StepTemp: 22}},
			},
		},
	}

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()
	_ = c.Configure(ctx, l)
	_ = c.StartFermentation(1)

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager}
	err := handler.Stop(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
}

func stopFermentationNotFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	err := handler.Stop(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func stopFermentationNotFermenting(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{
		ID: chamberID,
		CurrentBatch: &batch.Batch{
			Fermentation: batch.Fermentation{
				Steps: []batch.FermentationStep{{StepTemp: 22}},
			},
		},
	}

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	err := handler.Stop(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf("chamber '%s' is not fermenting", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func stopFermentationRespondError(t *testing.T) {
	t.Parallel()

	c := chamber.Chamber{
		ID: chamberID,
		CurrentBatch: &batch.Batch{
			Fermentation: batch.Fermentation{
				Steps: []batch.FermentationStep{{StepTemp: 22}},
			},
		},
	}

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()
	_ = c.Configure(ctx, l)
	_ = c.StartFermentation(1)

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAllChambers").Return([]chamber.Chamber{c}, nil)

	manager, _ := controller.NewChamberManager(ctx, repoMock, l)
	handler := &handlers.ChambersHandler{ChamberController: manager, Logger: l}
	// use new ctx to force error
	err := handler.Stop(context.Background(), w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}

func assetChamberListsAreEqual(t *testing.T, c1, c2 []chamber.Chamber) {
	t.Helper()
	assert.Equal(t, len(c1), len(c2))

	// source is a dictionary so order is not guaranteed
	sort.Slice(c1, func(i, j int) bool {
		return c1[i].ID < c1[j].ID
	})

	sort.Slice(c2, func(i, j int) bool {
		return c2[i].ID < c2[j].ID
	})

	for i := 0; i < len(c1); i++ {
		assertChambersAreEqual(t, c1[i], c2[i])
	}
}

func assertChambersAreEqual(t *testing.T, c1, c2 chamber.Chamber) {
	t.Helper()
	assert.Equal(t, c1.ID, c2.ID)
	assert.Equal(t, c1.Name, c2.Name)
	assert.Equal(t, c1.ChillerPin, c2.ChillerPin)
	assert.Equal(t, c1.HeaterPin, c2.HeaterPin)
	assert.Equal(t, c1.ChillerKp, c2.ChillerKp)
	assert.Equal(t, c1.ChillerKi, c2.ChillerKi)
	assert.Equal(t, c1.ChillerKd, c2.ChillerKd)
	assert.Equal(t, c1.HeaterKp, c2.HeaterKp)
	assert.Equal(t, c1.HeaterKi, c2.HeaterKi)
	assert.Equal(t, c1.HeaterKd, c2.HeaterKd)
	assert.Equal(t, c1.ModTime, c2.ModTime)
}
