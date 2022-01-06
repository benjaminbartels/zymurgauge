package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	brewfatherMock "github.com/benjaminbartels/zymurgauge/internal/test/mocks/brewfather"
	mocks "github.com/benjaminbartels/zymurgauge/internal/test/mocks/chamber"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	primaryStep = "Primary"
)

func getTestChamber() chamber.Chamber {
	return chamber.Chamber{
		ID: chamberID,
		DeviceConfigs: []chamber.DeviceConfig{
			{
				ID:    "1",
				Type:  "ds18b20",
				Roles: []string{"beerThermometer"},
			},
			{
				ID:    "2",
				Type:  "gpio",
				Roles: []string{"chiller"},
			},
			{
				ID:    "3",
				Type:  "gpio",
				Roles: []string{"heater"},
			},
		},
		CurrentBatch: &brewfather.Batch{
			Fermentation: brewfather.Fermentation{
				Steps: []brewfather.FermentationStep{
					{
						Type:     primaryStep,
						StepTemp: 22,
					},
				},
			},
		},
	}
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetAllChambers(t *testing.T) {
	t.Parallel()
	t.Run("getAllChambers", getAllChambers)
	t.Run("getAllChambersEmpty", getAllChambersEmpty)
	t.Run("getAllChamberOtherError", getAllChamberOtherError)
	t.Run("getAllChambersRespondError", getAllChambersRespondError)
}

func getAllChambers(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()
	chambers := []*chamber.Chamber{
		{ID: chamberID},
		{ID: "dd2610fe-95fc-45f3-8dd8-3051fb1bd4c1"},
	}

	controllerMock := &mocks.Controller{}
	controllerMock.On("GetAll").Return(chambers, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []*chamber.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)

	assertChamberListsAreEqual(t, chambers, result)
}

func getAllChambersEmpty(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()
	chambers := []*chamber.Chamber{}

	controllerMock := &mocks.Controller{}
	controllerMock.On("GetAll").Return(chambers, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []*chamber.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, chambers, result)
}

func getAllChamberOtherError(t *testing.T) {
	t.Parallel()

	c := &chamber.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(c)

	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("GetAll").Return(nil, errors.New("controllerMock error"))

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}

	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "get all chambers from"))
}

func getAllChambersRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("GetAll").Return([]*chamber.Chamber{}, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	// use new ctx to force error
	err := handler.GetAll(context.Background(), w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetChamber(t *testing.T) {
	t.Parallel()
	t.Run("getChamberFound", getChamberFound)
	t.Run("getChamberNotFoundError", getChamberNotFoundError)
	t.Run("getChamberOtherError", getChamberOtherError)
	t.Run("getChamberRespondError", getChamberRespondError)
}

func getChamberFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	c := &chamber.Chamber{ID: chamberID}
	controllerMock := &mocks.Controller{}
	controllerMock.On("Get", chamberID).Return(c, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}

	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)

	resp := w.Result()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := &chamber.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assertChambersAreEqual(t, c, result)
}

func getChamberNotFoundError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("Get", chamberID).Return(nil, nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}

	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func getChamberOtherError(t *testing.T) {
	t.Parallel()

	c := &chamber.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(c)

	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("Get", chamberID).Return(nil, errors.New("controllerMock error"))

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}

	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "get chamber "+chamberID+" from"))
}

func getChamberRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	c := &chamber.Chamber{ID: chamberID}
	controllerMock := &mocks.Controller{}
	controllerMock.On("Get", chamberID).Return(c, nil)

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
	t.Run("saveChamberInvalidConfigError", saveChamberInvalidConfigError)
	t.Run("saveChamberFermentingError", saveChamberFermentingError)
	t.Run("saveChamberOtherError", saveChamberOtherError)
	t.Run("saveChamberRespondError", saveChamberRespondError)
}

func saveChamber(t *testing.T) {
	t.Parallel()

	c := &chamber.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(c)

	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("Save", c).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}

	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := &chamber.Chamber{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
}

func saveChamberParseError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}

	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), parseErrorMsg)
}

func saveChamberInvalidConfigError(t *testing.T) {
	t.Parallel()

	c := getTestChamber()
	jsonBytes, _ := json.Marshal(c)

	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	chambers := []*chamber.Chamber{
		{ID: chamberID},
		{ID: "dd2610fe-95fc-45f3-8dd8-3051fb1bd4c1"},
	}

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(chambers, nil)

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(nil, errSomeError)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(nil, nil)

	serviceMock := &brewfatherMock.Service{}

	controller, err := chamber.NewManager(ctx, repoMock, configuratorMock, serviceMock, l)
	assert.NoError(t, err)

	handler := &handlers.ChambersHandler{ChamberController: controller, Logger: l}

	err = handler.Save(ctx, w, r, httprouter.Params{})

	assert.Contains(t, err.Error(), invalidConfigErrorMsg)

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func saveChamberFermentingError(t *testing.T) {
	t.Parallel()

	c := &chamber.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(c)

	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("Save", c).Return(chamber.ErrFermenting)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}

	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), fermentationInProgressMsg)

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func saveChamberOtherError(t *testing.T) {
	t.Parallel()

	c := &chamber.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(c)

	w, r, ctx := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("Save", c).Return(errors.New("controllerMock error"))

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}

	err := handler.Save(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "save chamber to"))
}

func saveChamberRespondError(t *testing.T) {
	t.Parallel()

	c := &chamber.Chamber{ID: chamberID}
	jsonBytes, _ := json.Marshal(c)

	w, r, _ := setupHandlerTest("", bytes.NewBuffer(jsonBytes))
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("Save", c).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}

	// use new ctx to force error
	err := handler.Save(context.Background(), w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestDeleteChamber(t *testing.T) {
	t.Parallel()
	t.Run("deleteChamber", deleteChamber)
	t.Run("deleteChamberNotFoundError", deleteChamberNotFoundError)
	t.Run("deleteChamberFermentingError", deleteChamberFermentingError)
	t.Run("deleteChamberOtherError", deleteChamberOtherError)
	t.Run("deleteChamberRespondError", deleteChamberRespondError)
}

func deleteChamber(t *testing.T) {
	t.Parallel()

	c := getTestChamber()

	w, r, ctx := setupHandlerTest("", nil)
	l, hook := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	serviceMock := &brewfatherMock.Service{}

	err := c.Configure(configuratorMock, serviceMock, l)
	assert.NoError(t, err)

	err = c.StartFermentation(ctx, "Primary")
	assert.NoError(t, err)

	controllerMock := &mocks.Controller{}
	controllerMock.On("Delete", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err = handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
	assert.Nil(t, hook.LastEntry())
}

func deleteChamberNotFoundError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("Delete", chamberID).Return(chamber.ErrNotFound)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)

	assert.Equal(t, http.StatusNotFound, reqErr.Status)
}

func deleteChamberFermentingError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("Delete", chamberID).Return(chamber.ErrFermenting)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fermentationInProgressMsg)

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)

	assert.Equal(t, http.StatusBadRequest, reqErr.Status)
}

func deleteChamberOtherError(t *testing.T) {
	t.Parallel()

	c := &chamber.Chamber{ID: chamberID}
	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	serviceMock := &brewfatherMock.Service{}

	err := c.Configure(configuratorMock, serviceMock, l)
	assert.NoError(t, err)

	controllerMock := &mocks.Controller{}
	controllerMock.On("Delete", chamberID).Return(errSomeError)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err = handler.Delete(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(controllerErrMsg, "delete chamber "+chamberID+" from"))
}

func deleteChamberRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("Delete", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	// use new ctx to force error
	err := handler.Delete(context.Background(), w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestStartFermentation(t *testing.T) {
	t.Parallel()
	t.Run("startFermentation", startFermentation)
	t.Run("startFermentationInvalidStepError", startFermentationInvalidStepError)
	t.Run("startFermentationNotFoundError", startFermentationNotFoundError)
	t.Run("startFermentationNoBatchError", startFermentationNoBatchError)
	t.Run("startFermentationOtherError", startFermentationOtherError)
	t.Run("startFermentationRespondError", startFermentationRespondError)
}

func startFermentation(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	c := getTestChamber()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	serviceMock := &brewfatherMock.Service{}

	err := c.Configure(configuratorMock, serviceMock, l)
	assert.NoError(t, err)

	w, r, ctx := setupHandlerTest("step="+primaryStep, nil)

	controllerMock := &mocks.Controller{}
	controllerMock.On("StartFermentation", chamberID, primaryStep).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock}
	err = handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
}

func startFermentationInvalidStepError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	c := getTestChamber()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	serviceMock := &brewfatherMock.Service{}

	err := c.Configure(configuratorMock, serviceMock, l)
	assert.NoError(t, err)

	step := "Secondary"
	w, r, ctx := setupHandlerTest("step="+step, nil)

	controllerMock := &mocks.Controller{}
	controllerMock.On("StartFermentation", chamberID, step).Return(chamber.ErrInvalidStep)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock}
	err = handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(invalidStepErrorMsg, step, chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func startFermentationNotFoundError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("step="+primaryStep, nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("StartFermentation", chamberID, primaryStep).Return(chamber.ErrNotFound)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func startFermentationNoBatchError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("step="+primaryStep, nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("StartFermentation", chamberID, primaryStep).Return(chamber.ErrNoCurrentBatch)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(noCurrentBatchErrorMsg, chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func startFermentationOtherError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("step="+primaryStep, nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("StartFermentation", chamberID, primaryStep).Return(errors.New("controllerMock error"))

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}

	err := handler.Start(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(startFermentationErrorMsg, chamberID))
}

func startFermentationRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest("step="+primaryStep, nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("StartFermentation", chamberID, primaryStep).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	// use new ctx to force error
	err := handler.Start(context.Background(), w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestStopFermentation(t *testing.T) {
	t.Parallel()
	t.Run("stopFermentation", stopFermentation)
	t.Run("stopFermentationNotFoundError", stopFermentationNotFoundError)
	t.Run("stopFermentationNotFermentingError", stopFermentationNotFermentingError)
	t.Run("stopFermentationOtherError", stopFermentationOtherError)
	t.Run("stopFermentationRespondError", stopFermentationRespondError)
}

func stopFermentation(t *testing.T) {
	t.Parallel()

	c := getTestChamber()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	serviceMock := &brewfatherMock.Service{}

	err := c.Configure(configuratorMock, serviceMock, l)
	assert.NoError(t, err)

	err = c.StartFermentation(ctx, "Primary")
	assert.NoError(t, err)

	controllerMock := &mocks.Controller{}
	controllerMock.On("StopFermentation", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock}
	err = handler.Stop(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.NoError(t, err)
}

func stopFermentationNotFoundError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("StopFermentation", chamberID).Return(chamber.ErrNotFound)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Stop(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "chamber", chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusNotFound)
}

func stopFermentationNotFermentingError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("StopFermentation", chamberID).Return(chamber.ErrNotFermenting)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	err := handler.Stop(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFermentingErrorMsg, chamberID))

	var reqErr *web.RequestError

	assert.ErrorAs(t, err, &reqErr)
	assert.Equal(t, reqErr.Status, http.StatusBadRequest)
}

func stopFermentationOtherError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	controllerMock := &mocks.Controller{}
	controllerMock.On("StopFermentation", chamberID).Return(errors.New("controllerMock error"))

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}

	err := handler.Stop(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(stopFermentationErrorMsg, chamberID))
}

func stopFermentationRespondError(t *testing.T) {
	t.Parallel()

	c := getTestChamber()

	w, r, ctx := setupHandlerTest("", nil)
	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	serviceMock := &brewfatherMock.Service{}

	err := c.Configure(configuratorMock, serviceMock, l)
	assert.NoError(t, err)

	err = c.StartFermentation(ctx, primaryStep)
	assert.NoError(t, err)

	controllerMock := &mocks.Controller{}
	controllerMock.On("StopFermentation", chamberID).Return(nil)

	handler := &handlers.ChambersHandler{ChamberController: controllerMock, Logger: l}
	// use new ctx to force error
	err = handler.Stop(context.Background(), w, r, httprouter.Params{httprouter.Param{Key: "id", Value: chamberID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}

func assertChamberListsAreEqual(t *testing.T, c1, c2 []*chamber.Chamber) {
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

func assertChambersAreEqual(t *testing.T, c1, c2 *chamber.Chamber) {
	t.Helper()
	assert.Equal(t, c1.ID, c2.ID)
	assert.Equal(t, c1.Name, c2.Name)
	assert.Equal(t, c1.ChillerKp, c2.ChillerKp)
	assert.Equal(t, c1.ChillerKi, c2.ChillerKi)
	assert.Equal(t, c1.ChillerKd, c2.ChillerKd)
	assert.Equal(t, c1.HeaterKp, c2.HeaterKp)
	assert.Equal(t, c1.HeaterKi, c2.HeaterKi)
	assert.Equal(t, c1.HeaterKd, c2.HeaterKd)
	assert.Equal(t, c1.ModTime, c2.ModTime)
}
