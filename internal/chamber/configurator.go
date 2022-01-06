//go:build !linux || !arm
// +build !linux !arm

package chamber

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
	"github.com/benjaminbartels/zymurgauge/internal/test/stubs"
)

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

var _ Configurator = (*DefaultConfigurator)(nil)

type DefaultConfigurator struct {
	TiltMonitor *tilt.Monitor
}

func (c *DefaultConfigurator) CreateDs18b20(thermometerID string) (device.Thermometer, error) {
	return &stubs.StubThermometer{ThermometerID: thermometerID}, nil
}

func (c *DefaultConfigurator) CreateTilt(color tilt.Color) (device.ThermometerAndHydrometer, error) {
	return &stubs.StubTilt{Color: color}, nil
}

func (c *DefaultConfigurator) CreateGPIOActuator(pin string) (device.Actuator, error) {
	return &stubs.StubGPIOActuator{Pin: pin}, nil
}
