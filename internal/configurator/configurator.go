//go:build !linux || !arm
// +build !linux !arm

package configurator

import "github.com/benjaminbartels/zymurgauge/internal/device/tilt"

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

// TODO: Force stubs from unit tests?

type Configurator struct{}

func (c Configurator) CreateDs18b20(thermometerID string) (*StubThermometer, error) {
	return &StubThermometer{thermometerID: thermometerID}, nil
}

func (c Configurator) CreateTilt(color tilt.Color) (*StubTilt, error) {
	return &StubTilt{color: color}, nil
}

func (c Configurator) CreateGPIOActuator(pin string) (*StubGPIOActuator, error) {
	return &StubGPIOActuator{pin: pin}, nil
}
