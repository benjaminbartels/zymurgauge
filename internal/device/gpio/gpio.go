package gpio

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/pkg/errors"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
)

// See: https://www.homebrewtalk.com/threads/craftbeerpi-raspberry-pi-software.569497/page-18
const frequency = 1 * physic.Hertz

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

func (a *Actuator) PWMOn(duty float64) error {
	d := gpio.Duty(float64(gpio.DutyMax) * duty)

	if err := a.pin.PWM(d, frequency); err != nil {
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
