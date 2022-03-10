package chamber_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

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
	chamberID1            = "96f58a65-03c0-49f3-83ca-ab751bbf3768"
	chamberID2            = "dd2610fe-95fc-45f3-8dd8-3051fb1bd4c1"
	chamberID3            = "82d328bf-a9a2-4bf9-adce-b87e0bd92141"
	repoErrMsg            = "could not %s repository"
	readingUpdateInterval = 100 * time.Millisecond
)

func createTestChambers() []*chamber.Chamber {
	chamber1 := &chamber.Chamber{
		ID:   chamberID1,
		Name: "ChamberWithCompleteConfigWithBatch",
		CurrentBatch: &brewfather.BatchDetail{
			Recipe: brewfather.Recipe{
				Fermentation: brewfather.Fermentation{
					Steps: []brewfather.FermentationStep{
						{
							Type:            "Primary",
							StepTemperature: 22,
						},
						{
							Type:            "Secondary",
							StepTemperature: 20,
						},
					},
				},
			},
		},
		DeviceConfig: chamber.DeviceConfig{
			ChillerGPIO:              "GPIO2",
			HeaterGPIO:               "GPIO3",
			BeerThermometerType:      "tilt",
			BeerThermometerID:        "orange",
			AuxiliaryThermometerType: "ds18b20",
			AuxiliaryThermometerID:   "28-0000071cbc72",
			ExternalThermometerType:  "ds18b20",
			ExternalThermometerID:    "28-000007158912",
			HydrometerType:           "tilt",
			HydrometerID:             "orange",
		},
	}
	chamber2 := &chamber.Chamber{
		ID:   chamberID2,
		Name: "ChamberWithMinimumConfigWithoutBatch",
		DeviceConfig: chamber.DeviceConfig{
			ChillerGPIO:         "GPIO5",
			HeaterGPIO:          "GPIO6",
			BeerThermometerType: "ds18b20",
			BeerThermometerID:   "28-0000071cbc72",
		},
	}
	chamber3 := &chamber.Chamber{
		ID:   chamberID3,
		Name: "ChamberWithMinimumConfigWithBatch",
		DeviceConfig: chamber.DeviceConfig{
			ChillerGPIO:         "GPIO5",
			HeaterGPIO:          "GPIO6",
			BeerThermometerType: "ds18b20",
			BeerThermometerID:   "28-0000071cbc72",
		},
		CurrentBatch: &brewfather.BatchDetail{
			Recipe: brewfather.Recipe{
				Fermentation: brewfather.Fermentation{
					Steps: []brewfather.FermentationStep{
						{
							Type:            "Primary",
							StepTemperature: 22,
						},
						{
							Type:            "Secondary",
							StepTemperature: 20,
						},
					},
				},
			},
		},
	}

	return []*chamber.Chamber{chamber1, chamber2, chamber3}
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
	m := &mocks.Metrics{}
	m.On("Gauge", mock.Anything, mock.Anything).Return()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(nil, errors.New("repoMock error"))

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil)

	manager, err := chamber.NewManager(context.Background(), repoMock, configuratorMock, serviceMock, l, m,
		readingUpdateInterval)
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "get all chambers from"))
	assert.Nil(t, manager)
}

func newManagerConfigureErrors(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	testChambers[0].DeviceConfig = chamber.DeviceConfig{
		BeerThermometerType: "badType",
		HydrometerType:      "badType",
	}

	l, _ := logtest.NewNullLogger()
	m := &mocks.Metrics{}
	m.On("Gauge", mock.Anything, mock.Anything).Return()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(testChambers, nil)

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil)

	manager, err := chamber.NewManager(context.Background(), repoMock, configuratorMock, serviceMock, l, m,
		readingUpdateInterval)
	assert.Contains(t, err.Error(), "could not configure temperature controllers")
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

	manager, _, _ := setupManagerTest(t, testChambers)

	result, err := manager.GetAll()
	assert.NoError(t, err)
	assertChamberListsAreEqual(t, testChambers, result)
}

func getAllChambersEmpty(t *testing.T) {
	t.Parallel()

	chambers := []*chamber.Chamber{}
	manager, _, _ := setupManagerTest(t, chambers)

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

	manager, _, _ := setupManagerTest(t, testChambers)

	result, err := manager.Get(chamberID1)
	assert.NoError(t, err)
	assertChambersAreEqual(t, testChambers[0], result)
}

func getChamberNotFoundError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _, _ := setupManagerTest(t, testChambers)

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

	manager, repoMock, _ := setupManagerTest(t, testChambers)
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

	manager, repoMock, _ := setupManagerTest(t, testChambers)
	repoMock.On("Save", testChambers[0]).Return(nil)

	err := manager.StartFermentation(chamberID1, "Primary")
	assert.NoError(t, err)

	err = manager.Save(testChambers[0])
	assert.ErrorIs(t, err, chamber.ErrFermenting)
}

func saveChamberRepoError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, repoMock, _ := setupManagerTest(t, testChambers)
	repoMock.On("Save", testChambers[0]).Return(errors.New("repoMock error"))

	err := manager.Save(testChambers[0])
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "save chamber to"))
}

func saveChamberConfigureError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, repoMock, _ := setupManagerTest(t, testChambers)
	repoMock.On("Save", testChambers[0]).Return(nil)

	c, err := manager.Get(chamberID1)
	assert.NoError(t, err)

	c.DeviceConfig = chamber.DeviceConfig{
		BeerThermometerType: "badType",
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

	manager, repoMock, _ := setupManagerTest(t, testChambers)
	repoMock.On("Delete", chamberID1).Return(nil)

	err := manager.Delete(chamberID1)
	assert.NoError(t, err)

	result, err := manager.Get(chamberID1)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func deleteChamberFermentingError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, repoMock, _ := setupManagerTest(t, testChambers)
	repoMock.On("Delete", chamberID1).Return(nil)

	err := manager.StartFermentation(chamberID1, "Primary")
	assert.NoError(t, err)

	err = manager.Delete(chamberID1)
	assert.ErrorIs(t, err, chamber.ErrFermenting)
}

func deleteChamberRepoError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, repoMock, _ := setupManagerTest(t, testChambers)
	repoMock.On("Delete", chamberID1).Return(errors.New("repoMock error"))

	err := manager.Delete(chamberID1)
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "delete chamber "+chamberID1+" from"))
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
	t.Run("startFermentationOtherDevicesAreNil", startFermentationOtherDevicesAreNil)
}

func startFermentation(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	var counter int

	doneCh := make(chan struct{}, 1)
	manager, _, metricsMock := setupManagerTest(t, testChambers)
	metricsMock.ExpectedCalls = nil // Erase previous On() functions
	metricsMock.On("Gauge", fmt.Sprintf("zymurgauge.%s.auxiliary_temperature,sensor_id=", testChambers[0].Name),
		mock.Anything).Return()
	metricsMock.On("Gauge", fmt.Sprintf("zymurgauge.%s.external_temperature,sensor_id=", testChambers[0].Name),
		mock.Anything).Return()
	metricsMock.On("Gauge", fmt.Sprintf("zymurgauge.%s.hydrometer_gravity,sensor_id=", testChambers[0].Name),
		mock.Anything).Return()
	metricsMock.On("Gauge", fmt.Sprintf("zymurgauge.%s.beer_temperature,sensor_id=", testChambers[0].Name),
		mock.Anything).Return().Run(
		func(args mock.Arguments) {
			counter++

			if counter == 2 {
				doneCh <- struct{}{}
			}
		})
	metricsMock.On(mock.Anything, mock.Anything).Return()

	err := manager.StartFermentation(chamberID1, "Primary")
	assert.NoError(t, err)

	select {
	case <-doneCh:
	case <-time.After(5 * time.Second):
		assert.Fail(t, "metrics.Gauge for beer temperature should have been called twice by now")
	}
}

func startFermentationNextStep(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _, _ := setupManagerTest(t, testChambers)
	err := manager.StartFermentation(chamberID1, "Primary")
	assert.NoError(t, err)

	err = manager.StartFermentation(chamberID1, "Secondary")
	assert.NoError(t, err)
}

func startFermentationNotFoundError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _, _ := setupManagerTest(t, testChambers)
	err := manager.StartFermentation("", "Primary")
	assert.ErrorIs(t, err, chamber.ErrNotFound)
}

func startFermentationNoCurrentBatchError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()
	testChambers[0].CurrentBatch = nil
	manager, _, _ := setupManagerTest(t, testChambers)

	err := manager.StartFermentation(chamberID1, "Primary")
	assert.ErrorIs(t, err, chamber.ErrNoCurrentBatch)
}

func startFermentationInvalidStepError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _, _ := setupManagerTest(t, testChambers)
	err := manager.StartFermentation(chamberID1, "BadStep")
	assert.ErrorIs(t, err, chamber.ErrInvalidStep)
}

func startFermentationTemperatureControllerLogError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	l, hook := logtest.NewNullLogger()
	m := &mocks.Metrics{}
	m.On("Gauge", mock.Anything, mock.Anything).Return()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(testChambers, nil)

	doneCh := make(chan struct{}, 1)

	thermometerMock := &mocks.ThermometerAndHydrometer{}
	thermometerMock.On("GetGravity").Return(0.0, nil)
	thermometerMock.On("GetID").Return("")
	thermometerMock.On("GetTemperature").Return(0.0, errors.New("thermometerMock error")).Run(
		func(args mock.Arguments) {
			doneCh <- struct{}{}
		})

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateTilt", mock.Anything).Return(thermometerMock, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil)

	manager, err := chamber.NewManager(context.Background(), repoMock, configuratorMock, serviceMock, l, m,
		readingUpdateInterval)
	assert.NoError(t, err)

	go func() {
		for {
			if logContains(hook.AllEntries(), logrus.ErrorLevel, "could not run temperature controller for chamber") {
				doneCh <- struct{}{}

				return
			}

			<-time.After(100 * time.Millisecond)
		}
	}()

	err = manager.StartFermentation(chamberID1, "Primary")
	assert.NoError(t, err)
	select {
	case <-doneCh:
	case <-time.After(5 * time.Second):
		assert.Fail(t, "log should contain expected value by now")
	}
}

func startFermentationOtherDevicesAreNil(t *testing.T) {
	t.Parallel()

	doneCh := make(chan struct{}, 1)

	testChambers := createTestChambers()
	l, _ := logtest.NewNullLogger()
	m := &mocks.Metrics{}
	m.On("Gauge", mock.Anything, mock.Anything).Return().Run(
		func(args mock.Arguments) {
			doneCh <- struct{}{}
		})

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(testChambers, nil)

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetID").Return("")
	thermometerMock.On("GetTemperature").Return(25.0, nil)

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateTilt", mock.Anything).Return(nil, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(thermometerMock, nil)

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil)
	manager, err := chamber.NewManager(context.Background(), repoMock, configuratorMock, serviceMock, l, m,
		readingUpdateInterval)
	assert.NoError(t, err)

	err = manager.StartFermentation(chamberID3, "Primary")
	assert.NoError(t, err)

	<-doneCh

	_, err = manager.Get(chamberID3)
	assert.NoError(t, err)
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

	manager, _, _ := setupManagerTest(t, testChambers)

	err := manager.StartFermentation(chamberID1, "Primary")
	assert.NoError(t, err)

	err = manager.StopFermentation(chamberID1)
	assert.NoError(t, err)
}

func stopFermentationNotFoundError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _, _ := setupManagerTest(t, testChambers)
	err := manager.StopFermentation("")
	assert.ErrorIs(t, err, chamber.ErrNotFound)
}

func stopFermentationNotFermentingError(t *testing.T) {
	t.Parallel()

	testChambers := createTestChambers()

	manager, _, _ := setupManagerTest(t, testChambers)

	err := manager.StopFermentation(chamberID1)
	assert.ErrorIs(t, err, chamber.ErrNotFermenting)
	assert.Equal(t, err.Error(), chamber.ErrNotFermenting.Error()) // For coverage
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
	chambers []*chamber.Chamber) (*chamber.Manager, *mocks.ChamberRepo, *mocks.Metrics) {
	t.Helper()

	l, _ := logtest.NewNullLogger()
	metricsMock := &mocks.Metrics{}
	metricsMock.On("Gauge", mock.Anything, mock.Anything).Return()

	repoMock := &mocks.ChamberRepo{}
	repoMock.On("GetAll").Return(chambers, nil)

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil)

	manager, err := chamber.NewManager(context.Background(), repoMock, configuratorMock, serviceMock, l, metricsMock,
		readingUpdateInterval)
	assert.NoError(t, err)

	return manager, repoMock, metricsMock
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
