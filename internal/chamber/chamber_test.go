package chamber_test

import (
	"fmt"
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	brewfatherMocks "github.com/benjaminbartels/zymurgauge/internal/test/mocks/brewfather"
	mocks "github.com/benjaminbartels/zymurgauge/internal/test/mocks/chamber"
	"github.com/pkg/errors"
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

	err := c[0].Configure(configuratorMock, serviceMock, l)
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

	err := c[0].Configure(configuratorMock, serviceMock, l) // element 0 has ds18b20

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

	err := c[1].Configure(configuratorMock, serviceMock, l) // element 1 has tilt

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

	err := c[0].Configure(configuratorMock, serviceMock, l)

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

	err := c[0].Configure(configuratorMock, serviceMock, l)

	var cfgErr *chamber.InvalidConfigurationError

	assert.ErrorAs(t, err, &cfgErr)
	assert.Contains(t, cfgErr.Problems()[0].Error(), fmt.Sprintf(roleErrMsg, badRole))
}
