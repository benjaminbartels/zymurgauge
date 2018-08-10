package gpio

import (
	"github.com/stianeikeland/go-rpio"
)

// Actuator is a device that turns on and off by manipulating a GPIO pin
type Actuator struct {
	pin rpio.Pin
}

// NewActuator creates a new Actuator assigned to the given pinID
func NewActuator(pinID uint8) *Actuator {
	return &Actuator{pin: rpio.Pin(pinID)}
}

// On turns the Actuator on by setting the pin to high
func (a *Actuator) On() error {
	return a.pin.High()
}

// Off turns the Actuator on by setting the pin to high
func (a *Actuator) Off() error {
	return a.pin.Low()
}

// IsOn returns true if the Actuator is on
func (a *Actuator) IsOn() bool {
	return a.pin.Read() == rpio.High
}
