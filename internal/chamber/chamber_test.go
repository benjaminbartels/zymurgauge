package chamber_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/benjaminbartels/zymurgauge/internal/test/stubs"
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
)

//nolint:paralleltest // False positives with r.Run not in a loop
func TestConfigure(t *testing.T) {
	t.Parallel()
	t.Run("configure", configure)
	t.Run("configureDs18b20Error", configureDs18b20Error)
	t.Run("configureTiltError", configureTiltError)
	t.Run("configureGPIOError", configureGPIOError)
}

const (
	ds18b20ID = "28-0000071cbc72"
	tiltColor = "orange"
	gpio2     = "GPIO2"
	gpio3     = "GPIO3"
)

func configure(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	m := &mocks.Metrics{}
	m.On("Gauge", mock.Anything, mock.Anything).Return()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil)

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, l, m, readingUpdateInterval)
	assert.NoError(t, err)
}

func configureDs18b20Error(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	m := &mocks.Metrics{}
	m.On("Gauge", mock.Anything, mock.Anything).Return()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(nil, errors.New("configuratorMock error"))
	configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil)

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, l, m, readingUpdateInterval) // element 1 has ds18b20

	var cfgErr *chamber.InvalidConfigurationError

	assert.ErrorAs(t, err, &cfgErr)
	assert.Contains(t, cfgErr.Problems()[0].Error(), fmt.Sprintf(ds18b20ErrMsg, ds18b20ID))
}

func configureTiltError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	m := &mocks.Metrics{}
	m.On("Gauge", mock.Anything, mock.Anything).Return()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(nil, errors.New("configuratorMock error"))
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	c := createTestChambers()

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil)

	err := c[0].Configure(configuratorMock, serviceMock, l, m, readingUpdateInterval) // element 0 has tilt

	var cfgErr *chamber.InvalidConfigurationError

	assert.ErrorAs(t, err, &cfgErr)
	assert.Contains(t, cfgErr.Problems()[0].Error(), fmt.Sprintf(tiltErrMsg, tiltColor))
}

func configureGPIOError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	m := &mocks.Metrics{}
	m.On("Gauge", mock.Anything, mock.Anything).Return()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", gpio2).Return(nil, errors.New("configuratorMock error"))
	configuratorMock.On("CreateGPIOActuator", gpio3).Return(nil, errors.New("configuratorMock error"))

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil)

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, l, m, readingUpdateInterval)

	var cfgErr *chamber.InvalidConfigurationError

	assert.ErrorAs(t, err, &cfgErr)
	assert.Contains(t, cfgErr.Problems()[0].Error(), fmt.Sprintf(gpioErrMsg, gpio2))
}

//nolint:paralleltest // False positives with r.Run not in a loop
func TestLogging(t *testing.T) {
	t.Parallel()
	t.Run("log", log)
	t.Run("logServiceErrors", logServiceErrors)
	t.Run("logDeviceErrors", logDeviceErrors)
}

func log(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	m := &mocks.Metrics{}
	m.On("Gauge", mock.Anything, mock.Anything).Return()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	doneCh := make(chan struct{}, 1)

	expected := brewfather.LogEntry{
		Beer:                 "Pale Ale",
		DeviceName:           "ChamberWithCompleteConfigWithBatch",
		BeerTemperature:      "25.000000",
		AuxiliaryTemperature: "25.000000",
		ExternalTemperature:  "25.000000",
		Gravity:              "0.950000",
		TemperatureUnit:      "C",
		GravityUnit:          "G",
	}

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			assert.Equal(t, expected, args[1])
			doneCh <- struct{}{}
		})

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, l, m, readingUpdateInterval)
	assert.NoError(t, err)

	err = c[0].StartFermentation(context.Background(), "Primary")
	assert.NoError(t, err)

	<-doneCh
}

func logServiceErrors(t *testing.T) {
	t.Parallel()

	l, hook := logtest.NewNullLogger()
	m := &mocks.Metrics{}
	m.On("Gauge", mock.Anything, mock.Anything).Return()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&stubs.Thermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&stubs.Tilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	doneCh := make(chan struct{}, 1)

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(errors.New("serviceMock error")).Run(
		func(args mock.Arguments) {
			defer func() { doneCh <- struct{}{} }()
		})

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, l, m, readingUpdateInterval)
	assert.NoError(t, err)

	err = c[0].StartFermentation(context.Background(), "Primary")
	assert.NoError(t, err)

	<-doneCh

	assert.True(t, logContains(hook.AllEntries(), logrus.ErrorLevel, "Unable to send readings to Brewfather"))
}

func logDeviceErrors(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()
	m := &mocks.Metrics{}
	m.On("Gauge", mock.Anything, mock.Anything).Return()

	tiltMock := &mocks.ThermometerAndHydrometer{}
	tiltMock.On("GetTemperature").Return(0.0, errors.New("tiltMock error"))
	tiltMock.On("GetGravity").Return(0.0, errors.New("tiltMock error"))

	thermometerMock := &mocks.Thermometer{}
	thermometerMock.On("GetTemperature").Return(0.0, errors.New("thermometerMock error"))

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(thermometerMock, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(tiltMock, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&stubs.Actuator{}, nil)

	doneCh := make(chan struct{}, 1)
	ctx, stop := context.WithCancel(context.Background())

	expected := brewfather.LogEntry{
		Beer:            "Pale Ale",
		DeviceName:      "ChamberWithCompleteConfigWithBatch",
		TemperatureUnit: "C",
		GravityUnit:     "G",
	}

	serviceMock := &mocks.Service{}
	serviceMock.On("Log", mock.Anything, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			assert.Equal(t, expected, args[1])

			doneCh <- struct{}{}
			stop()
		})

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, serviceMock, l, m, readingUpdateInterval)
	assert.NoError(t, err)

	err = c[0].StartFermentation(ctx, "Primary")
	assert.NoError(t, err)

	<-doneCh
}
