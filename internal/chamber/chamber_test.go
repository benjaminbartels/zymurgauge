package chamber_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	brewfatherMocks "github.com/benjaminbartels/zymurgauge/internal/test/mocks/brewfather"
	mocks "github.com/benjaminbartels/zymurgauge/internal/test/mocks/chamber"
	deviceMocks "github.com/benjaminbartels/zymurgauge/internal/test/mocks/device"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	ds18b20ErrMsg = "could not create new Ds18b20 %s"
	tiltErrMsg    = "could not create new %s Tilt"
	gpioErrMsg    = "could not create new GPIO %s"
	roleErrMsg    = "invalid device role '%s'"
)

//nolint: paralleltest // False positives with r.Run not in a loop
func TestConfigure(t *testing.T) {
	t.Parallel()
	t.Run("configure", configure)
	t.Run("configureDs18b20Error", configureDs18b20Error)
	t.Run("configureTiltError", configureTiltError)
	t.Run("configureGPIOError", configureGPIOError)
	t.Run("configureInvalidRoleError", configureInvalidRoleError)
}

const (
	ds18b20ID = "28-0000071cbc72"
	tiltColor = "orange"
	gpio2     = "GPIO2"
	gpio3     = "GPIO3"
	badRole   = "badRole"
)

func configure(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	serviceMock := &brewfatherMocks.Service{}

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, false, l)
	assert.NoError(t, err)
}

func configureDs18b20Error(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(nil, errors.New("configuratorMock error"))
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	serviceMock := &brewfatherMocks.Service{}

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, false, l) // element 1 has ds18b20

	var cfgErr *chamber.InvalidConfigurationError

	assert.ErrorAs(t, err, &cfgErr)
	assert.Contains(t, cfgErr.Problems()[0].Error(), fmt.Sprintf(ds18b20ErrMsg, ds18b20ID))
}

func configureTiltError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(nil, errors.New("configuratorMock error"))
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	c := createTestChambers()

	serviceMock := &brewfatherMocks.Service{}

	err := c[0].Configure(configuratorMock, serviceMock, false, l) // element 0 has tilt

	var cfgErr *chamber.InvalidConfigurationError

	assert.ErrorAs(t, err, &cfgErr)
	assert.Contains(t, cfgErr.Problems()[0].Error(), fmt.Sprintf(tiltErrMsg, tiltColor))
}

func configureGPIOError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", gpio2).Return(nil, errors.New("configuratorMock error"))
	configuratorMock.On("CreateGPIOActuator", gpio3).Return(nil, errors.New("configuratorMock error"))

	serviceMock := &brewfatherMocks.Service{}

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, false, l)

	var cfgErr *chamber.InvalidConfigurationError

	assert.ErrorAs(t, err, &cfgErr)
	assert.Contains(t, cfgErr.Problems()[0].Error(), fmt.Sprintf(gpioErrMsg, gpio2))
}

func configureInvalidRoleError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	serviceMock := &brewfatherMocks.Service{}

	c := createTestChambers()

	c[0].DeviceConfigs[0].Roles[0] = "badRole"

	err := c[0].Configure(configuratorMock, serviceMock, false, l)

	var cfgErr *chamber.InvalidConfigurationError

	assert.ErrorAs(t, err, &cfgErr)
	assert.Contains(t, cfgErr.Problems()[0].Error(), fmt.Sprintf(roleErrMsg, badRole))
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestLogging(t *testing.T) {
	t.Parallel()
	t.Run("log", log)
	t.Run("logServiceErrors", logServiceErrors)
	t.Run("logDeviceErrors", logDeviceErrors)
}

func log(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	doneCh := make(chan struct{}, 1)

	expected := brewfather.LogEntry{
		DeviceName:           "Chamber1",
		BeerTemperature:      "25.000000",
		AuxiliaryTemperature: "25.000000",
		ExternalTemperature:  "25.000000",
		Gravity:              "0.950000",
		TemperatureUnit:      "C",
		GravityUnit:          "G",
	}

	serviceMock := &brewfatherMocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			assert.Equal(t, expected, args[1])
			doneCh <- struct{}{}
		})

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, true, l)
	assert.NoError(t, err)

	err = c[0].StartFermentation(context.Background(), "Primary")
	assert.NoError(t, err)

	<-doneCh
}

func logServiceErrors(t *testing.T) {
	t.Parallel()

	l, hook := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	doneCh := make(chan struct{}, 1)

	serviceMock := &brewfatherMocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(errors.New("thermometerMock error")).Run(
		func(args mock.Arguments) {
			doneCh <- struct{}{}
		})

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, true, l)
	assert.NoError(t, err)

	err = c[0].StartFermentation(context.Background(), "Primary")
	assert.NoError(t, err)

	<-doneCh

	assert.True(t, logContains(hook.AllEntries(), logrus.ErrorLevel, "could not log tilt data"))
}

func logDeviceErrors(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	tiltMock := &deviceMocks.ThermometerAndHydrometer{}
	tiltMock.On("GetTemperature").Return(0.0, errors.New("tiltMock error"))
	tiltMock.On("GetGravity").Return(0.0, errors.New("tiltMock error"))

	thermometerMock := &deviceMocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(0.0, errors.New("thermometerMock error"))

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(thermometerMock, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(tiltMock, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	doneCh := make(chan struct{}, 1)
	ctx, stop := context.WithCancel(context.Background())

	expected := brewfather.LogEntry{
		DeviceName:      "Chamber1",
		TemperatureUnit: "C",
		GravityUnit:     "G",
	}

	serviceMock := &brewfatherMocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			assert.Equal(t, expected, args[1])

			doneCh <- struct{}{}
			stop()
		})

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, true, l)
	assert.NoError(t, err)

	err = c[0].StartFermentation(ctx, "Primary")
	assert.NoError(t, err)

	<-doneCh
}
