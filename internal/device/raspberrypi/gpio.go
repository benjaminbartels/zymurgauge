package raspberrypi

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/pkg/errors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

var _ device.Actuator = (*GPIOActuator)(nil)

type GPIOActuator struct {
	pin gpio.PinIO
}

func NewGPIOActuator(pinID string) (*GPIOActuator, error) {
	pin := gpioreg.ByName(pinID)
	if pin == nil {
		return nil, errors.Errorf("Could not open %s", pinID)
	}

	return &GPIOActuator{pin: pin}, nil
}

func (a *GPIOActuator) On() error {
	if err := a.pin.Out(gpio.High); err != nil {
		return errors.Wrapf(err, "could not set pin %s to high", a.pin.Name())
	}

	return nil
}

func (a *GPIOActuator) Off() error {
	if err := a.pin.Out(gpio.Low); err != nil {
		return errors.Wrapf(err, "could not set pin %s to low", a.pin.Name())
	}

	return nil
}
