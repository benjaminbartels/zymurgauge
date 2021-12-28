package gpio

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/pkg/errors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

var _ device.Actuator = (*Actuator)(nil)

type Actuator struct {
	pin gpio.PinIO
}

func NewGPIOActuator(pinID string) (*Actuator, error) {
	pin := gpioreg.ByName(pinID)
	if pin == nil {
		return nil, errors.Errorf("Could not open %s", pinID)
	}

	return &Actuator{pin: pin}, nil
}

func (a *Actuator) On() error {
	if err := a.pin.Out(gpio.High); err != nil {
		return errors.Wrapf(err, "could not set pin %s to high", a.pin.Name())
	}

	return nil
}

func (a *Actuator) Off() error {
	if err := a.pin.Out(gpio.Low); err != nil {
		return errors.Wrapf(err, "could not set pin %s to low", a.pin.Name())
	}

	return nil
}
