package chamber

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/gpio"
	"github.com/benjaminbartels/zymurgauge/internal/device/onewire"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
	"github.com/pkg/errors"
)

var _ Configurator = (*DefaultConfigurator)(nil)

type DefaultConfigurator struct {
	TiltMonitor *tilt.Monitor
}

func (c *DefaultConfigurator) CreateDs18b20(thermometerID string) (device.Thermometer, error) {
	ds18b20, err := onewire.NewDs18b20(thermometerID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create new Ds18b20 thermometer %s", thermometerID)
	}

	return ds18b20, nil
}

func (c *DefaultConfigurator) CreateTilt(color tilt.Color) (device.ThermometerAndHydrometer, error) {
	tilt, err := c.TiltMonitor.GetTilt(color)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get %s tilt", color)
	}

	return tilt, nil
}

func (c *DefaultConfigurator) CreateGPIOActuator(pin string) (device.Actuator, error) {
	actuator, err := gpio.NewGPIOActuator(pin)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create new raspberry pi gpio actuator for pin %s", pin)
	}

	return actuator, nil
}
