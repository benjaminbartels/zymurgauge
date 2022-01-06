package chamber_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/benjaminbartels/zymurgauge/internal/test/stubs"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	chamberID  = "96f58a65-03c0-49f3-83ca-ab751bbf3768"
	repoErrMsg = "could not %s repository"
)

func createTestChambers() []*chamber.Chamber {
	chamber1 := chamber.Chamber{
		ID:   chamberID,
		Name: "Chamber1",
		CurrentBatch: &brewfather.Batch{
			Fermentation: brewfather.Fermentation{
				Steps: []brewfather.FermentationStep{
					{
						Type:     "Primary",
						StepTemp: 22,
					},
					{
						Type:     "Secondary",
						StepTemp: 20,
					},
				},
			},
		},
		DeviceConfigs: []chamber.DeviceConfig{
			{ID: "orange", Type: "tilt", Roles: []string{"beerThermometer", "hydrometer"}},
			{ID: "28-0000071cbc72", Type: "ds18b20", Roles: []string{"auxiliaryThermometer"}},
			{ID: "28-000007158912", Type: "ds18b20", Roles: []string{"externalThermometer"}},
			{ID: "GPIO2", Type: "gpio", Roles: []string{"chiller"}},
			{ID: "GPIO3", Type: "gpio", Roles: []string{"heater"}},
		},
	}
	chamber2 := chamber.Chamber{
		ID:   "dd2610fe-95fc-45f3-8dd8-3051fb1bd4c1",
		Name: "Chamber2",
		DeviceConfigs: []chamber.DeviceConfig{
			{ID: "28-0000071cbc72", Type: "ds18b20", Roles: []string{"beerThermometer"}},
			{ID: "GPIO5", Type: "gpio", Roles: []string{"chiller"}},
			{ID: "GPIO6", Type: "gpio", Roles: []string{"heater"}},
		},
	}

	return []*chamber.Chamber{&chamber1, &chamber2}
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestNewManager(t *testing.T) {
	t.Parallel()
	t.Run("newManagerGetAllError", newManagerGetAllError)
	t.Run("newManagerConfigureErrors", newManagerConfigureErrors)
}

func newManagerGetAllError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(nil, errors.New("repoMock error"))

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	serviceMock := &mocks.Service{}

	manager, err := chamber.NewManager(context.Background(), repoMock, configuratorMock, serviceMock, false, l)
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "get all chambers from"))
	assert.Nil(t, manager)
}

func newManagerConfigureErrors(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	testChambers[0].DeviceConfigs = []chamber.DeviceConfig{
		{ID: "1", Type: "badType", Roles: []string{}},
	}

	l, _ := logtest.NewNullLogger()
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(testChambers, nil)

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	serviceMock := &mocks.Service{}

	manager, err := chamber.NewManager(context.Background(), repoMock, configuratorMock, serviceMock, false, l)
	assert.Contains(t, err.Error(), "could not configure all temperature controllers")
	assert.NotNil(t, manager)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetAllChambers(t *testing.T) {
	t.Parallel()
	t.Run("getAllChambers", getAllChambers)
	t.Run("getAllChambersEmpty", getAllChambersEmpty)
}

func getAllChambers(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _ := setupManagerTest(t, testChambers)

	result, err := manager.GetAll()
	assert.NoError(t, err)
	assertChamberListsAreEqual(t, testChambers, result)
}

func getAllChambersEmpty(t *testing.T) {
	t.Parallel()

	chambers := []*chamber.Chamber{}
	manager, _ := setupManagerTest(t, chambers)

	result, err := manager.GetAll()
	assert.NoError(t, err)
	assertChamberListsAreEqual(t, chambers, result)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetChamber(t *testing.T) {
	t.Parallel()
	t.Run("getChamberFound", getChamberFound)
	t.Run("getChamberNotFoundError", getChamberNotFoundError)
}

func getChamberFound(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _ := setupManagerTest(t, testChambers)

	result, err := manager.Get(chamberID)
	assert.NoError(t, err)
	assertChambersAreEqual(t, testChambers[0], result)
}

func getChamberNotFoundError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _ := setupManagerTest(t, testChambers)

	result, err := manager.Get("")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestSaveChamber(t *testing.T) {
	t.Parallel()
	t.Run("saveChamber", saveChamber)
	t.Run("saveChamberFermentingError", saveChamberFermentingError)
	t.Run("saveChamberRepoError", saveChamberRepoError)
	t.Run("saveChamberConfigureError", saveChamberConfigureError)
}

func saveChamber(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	newChamberID := "a605c238-e6b6-4ebc-b997-965d29d47060"

	newChamber := &chamber.Chamber{
		ID: newChamberID,
	}

	manager, repoMock := setupManagerTest(t, testChambers)
	repoMock.On("Save", newChamber).Return(nil)

	err := manager.Save(newChamber)
	assert.NoError(t, err)

	result, err := manager.Get(newChamberID)
	assert.NoError(t, err)
	assertChambersAreEqual(t, newChamber, result)
}

func saveChamberFermentingError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, repoMock := setupManagerTest(t, testChambers)
	repoMock.On("Save", testChambers[0]).Return(nil)

	err := manager.StartFermentation(chamberID, "Primary")
	assert.NoError(t, err)

	err = manager.Save(testChambers[0])
	assert.ErrorIs(t, err, chamber.ErrFermenting)
}

func saveChamberRepoError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, repoMock := setupManagerTest(t, testChambers)
	repoMock.On("Save", testChambers[0]).Return(errors.New("repoMock error"))

	err := manager.Save(testChambers[0])
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "save chamber to"))
}

func saveChamberConfigureError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, repoMock := setupManagerTest(t, testChambers)
	repoMock.On("Save", testChambers[0]).Return(nil)

	c, err := manager.Get(chamberID)
	assert.NoError(t, err)

	c.DeviceConfigs = []chamber.DeviceConfig{
		{ID: "1", Type: "badType", Roles: []string{}},
	}

	err = manager.Save(c)

	var cfgErr *chamber.InvalidConfigurationError

	assert.ErrorAs(t, err, &cfgErr)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestDeleteChamber(t *testing.T) {
	t.Parallel()
	t.Run("deleteChamber", deleteChamber)
	t.Run("deleteChamberFermentingError", deleteChamberFermentingError)
	t.Run("deleteChamberRepoError", deleteChamberRepoError)
}

func deleteChamber(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, repoMock := setupManagerTest(t, testChambers)
	repoMock.On("Delete", chamberID).Return(nil)

	err := manager.Delete(chamberID)
	assert.NoError(t, err)

	result, err := manager.Get(chamberID)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func deleteChamberFermentingError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, repoMock := setupManagerTest(t, testChambers)
	repoMock.On("Delete", chamberID).Return(nil)

	err := manager.StartFermentation(chamberID, "Primary")
	assert.NoError(t, err)

	err = manager.Delete(chamberID)
	assert.ErrorIs(t, err, chamber.ErrFermenting)
}

func deleteChamberRepoError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, repoMock := setupManagerTest(t, testChambers)
	repoMock.On("Delete", chamberID).Return(errors.New("repoMock error"))

	err := manager.Delete(chamberID)
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "delete chamber "+chamberID+" from"))
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestStartFermentation(t *testing.T) {
	t.Parallel()
	t.Run("startFermentation", startFermentation)
	t.Run("startFermentationNextStep", startFermentationNextStep)
	t.Run("startFermentationNotFoundError", startFermentationNotFoundError)
	t.Run("startFermentationNoCurrentBatchError", startFermentationNoCurrentBatchError)
	t.Run("startFermentationInvalidStepError", startFermentationInvalidStepError)
	t.Run("startFermentationTemperatureControllerLogError", startFermentationTemperatureControllerLogError)
}

func startFermentation(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _ := setupManagerTest(t, testChambers)
	err := manager.StartFermentation(chamberID, "Primary")
	assert.NoError(t, err)
}

func startFermentationNextStep(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _ := setupManagerTest(t, testChambers)
	err := manager.StartFermentation(chamberID, "Primary")
	assert.NoError(t, err)

	err = manager.StartFermentation(chamberID, "Secondary")
	assert.NoError(t, err)
}

func startFermentationNotFoundError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _ := setupManagerTest(t, testChambers)
	err := manager.StartFermentation("", "Primary")
	assert.ErrorIs(t, err, chamber.ErrNotFound)
}

func startFermentationNoCurrentBatchError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()
	testChambers[0].CurrentBatch = nil
	manager, _ := setupManagerTest(t, testChambers)

	err := manager.StartFermentation(chamberID, "Primary")
	assert.ErrorIs(t, err, chamber.ErrNoCurrentBatch)
}

func startFermentationInvalidStepError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _ := setupManagerTest(t, testChambers)
	err := manager.StartFermentation(chamberID, "BadStep")
	assert.ErrorIs(t, err, chamber.ErrInvalidStep)
}

func startFermentationTemperatureControllerLogError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	l, hook := logtest.NewNullLogger()
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(testChambers, nil)

	doneCh := make(chan struct{}, 1)

	thermometerMock := &mocks.ThermometerAndHydrometer{}
	thermometerMock.On("GetTemperature").Return(0.0, errors.New("thermometerMock error")).Run(
		func(args mock.Arguments) {
			doneCh <- struct{}{}
		})

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateTilt", mock.Anything).Return(thermometerMock, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)

	serviceMock := &mocks.Service{}

	manager, err := chamber.NewManager(context.Background(), repoMock, configuratorMock, serviceMock, false, l)
	assert.NoError(t, err)

	err = manager.StartFermentation(chamberID, "Primary")
	assert.NoError(t, err)

	<-doneCh

	logContains(hook.AllEntries(), logrus.ErrorLevel, "could not run temperature controller for chamber")
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestStopFermentation(t *testing.T) {
	t.Parallel()
	t.Run("stopFermentation", stopFermentation)
	t.Run("stopFermentationNotFoundError", stopFermentationNotFoundError)
	t.Run("stopFermentationNotFermentingError", stopFermentationNotFermentingError)
}

func stopFermentation(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _ := setupManagerTest(t, testChambers)

	err := manager.StartFermentation(chamberID, "Primary")
	assert.NoError(t, err)

	err = manager.StopFermentation(chamberID)
	assert.NoError(t, err)
}

func stopFermentationNotFoundError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _ := setupManagerTest(t, testChambers)
	err := manager.StopFermentation("")
	assert.ErrorIs(t, err, chamber.ErrNotFound)
}

func stopFermentationNotFermentingError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _ := setupManagerTest(t, testChambers)

	err := manager.StopFermentation(chamberID)
	assert.ErrorIs(t, err, chamber.ErrNotFermenting)
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

func setupManagerTest(t *testing.T,
	chambers []*chamber.Chamber) (*chamber.Manager, *mocks.ChamberRepo) {
	t.Helper()

	l, _ := logtest.NewNullLogger()
	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(chambers, nil)

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	serviceMock := &mocks.Service{}

	manager, err := chamber.NewManager(context.Background(), repoMock, configuratorMock, serviceMock, false, l)
	assert.NoError(t, err)

	return manager, repoMock
}

func logContains(logs []*logrus.Entry, level logrus.Level, substr string) bool {
	found := false

	for _, v := range logs {
		if strings.Contains(v.Message, substr) && v.Level == level {
			found = true
		}
	}

	return found
}
