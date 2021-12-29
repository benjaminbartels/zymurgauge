package configurator

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/gpio"
	"github.com/benjaminbartels/zymurgauge/internal/device/onewire"
	"github.com/pkg/errors"
)

type Configurator struct{}

func (c Configurator) CreateDs18b20(thermometerID string) (*onewire.Ds18b20, error) {
	ds18b20, err := onewire.NewDs18b20(thermometerID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create new Ds18b20 thermometer %s", thermometerID)
	}

	return ds18b20, nil
}

func (c Configurator) CreateTilt(thermometerID string) (device.Thermometer, error) {
	ds18b20, err := onewire.NewDs18b20(thermometerID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create new Ds18b20 thermometer %s", thermometerID)
	}

	return ds18b20, nil
}

func (c Configurator) CreateGPIOActuator(pin string) (device.Actuator, error) {
	actuator, err := gpio.NewGPIOActuator(pin)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create new raspberry pi gpio actuator for pin %s", pin)
	}

	return actuator, nil
}
