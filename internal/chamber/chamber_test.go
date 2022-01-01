package chamber_test

import (
	"fmt"
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/chamber"
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
	gpio02    = "GPIO2"
)

func configure(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, l)
	assert.NoError(t, err)
}

func configureDs18b20Error(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(nil, errors.New("repoMock error"))
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, l) // element 0 has ds18b20
	assert.Contains(t, err.Error(), fmt.Sprintf(ds18b20ErrMsg, ds18b20ID))
}

func configureTiltError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(nil, errors.New("repoMock error"))
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	c := createTestChambers()

	err := c[1].Configure(configuratorMock, l) // element 1 has tilt
	assert.Contains(t, err.Error(), fmt.Sprintf(tiltErrMsg, tiltColor))
}

func configureGPIOError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", gpio02).Return(nil, errors.New("repoMock error"))

	c := createTestChambers()

	err := c[0].Configure(configuratorMock, l)
	assert.Contains(t, err.Error(), fmt.Sprintf(gpioErrMsg, gpio02))
}

func configureInvalidRoleError(t *testing.T) {
	t.Parallel()

	l, _ := logtest.NewNullLogger()

	configuratorMock := &mocks.Configurator{}
	configuratorMock.On("CreateDs18b20", mock.Anything).Return(&chamber.StubThermometer{}, nil)
	configuratorMock.On("CreateTilt", mock.Anything).Return(&chamber.StubTilt{}, nil)
	configuratorMock.On("CreateGPIOActuator", mock.Anything).Return(&chamber.StubGPIOActuator{}, nil)

	c := createTestChambers()

	c[0].DeviceConfigs[0].Roles[0] = "badRole"

	err := c[0].Configure(configuratorMock, l)
	assert.ErrorIs(t, err, chamber.ErrInvalidDeviceRole)
}
