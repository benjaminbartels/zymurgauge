package internal

import (
	"github.com/pkg/errors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type GPIOActuator struct {
	pin gpio.PinIO
}

func NewGPIOActuator(pinID string) (*GPIOActuator, error) {

	pin := gpioreg.ByName(pinID)
	if pin == nil {
		return nil, errors.Errorf("Could not open %s", pin)
	}

	return &GPIOActuator{pin: pin}, nil
}

func (a *GPIOActuator) On() error {
	return a.pin.Out(gpio.High)
}

func (a *GPIOActuator) Off() error {
	return a.pin.Out(gpio.Low)
}

func (a *GPIOActuator) IsOn() bool {
	return a.pin.Read() == gpio.High
}
