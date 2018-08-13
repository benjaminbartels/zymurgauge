package gpio

import (
	"github.com/pkg/errors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

// ToDo: Make constants for pinID

// Actuator is a device that turns on and off by manipulating a GPIO pin
type Actuator struct {
	pin gpio.PinIO
}

// NewActuator creates a new Actuator assigned to the give pinID
func NewActuator(pinID string) (*Actuator, error) {

	pin := gpioreg.ByName(pinID)
	if pin == nil {
		return nil, errors.Errorf("Could not create new Actuator. Could not open %s", pin)
	}

	return &Actuator{pin: pin}, nil
}

// On turns the Actuator on by setting the pin to high
func (a *Actuator) On() error {
	return a.pin.Out(gpio.High)
}

// Off turns the Actuator on by setting the pin to high
func (a *Actuator) Off() error {
	return a.pin.Out(gpio.Low)
}

// IsOn returns true if the Actuator is on
func (a *Actuator) IsOn() bool {
	return a.pin.Read() == gpio.High
}
