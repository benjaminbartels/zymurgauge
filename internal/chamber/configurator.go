//go:build !linux || !arm
// +build !linux !arm

package chamber

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
)

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

type DefaultConfigurator struct {
	TiltMonitor tilt.Monitor
}

func (c *DefaultConfigurator) CreateDs18b20(thermometerID string) (device.Thermometer, error) {
	return &StubThermometer{thermometerID: thermometerID}, nil
}

func (c *DefaultConfigurator) CreateTilt(color tilt.Color) (device.ThermometerAndHydrometer, error) {
	return &StubTilt{color: color}, nil
}

func (c *DefaultConfigurator) CreateGPIOActuator(pin string) (device.Actuator, error) {
	return &StubGPIOActuator{pin: pin}, nil
}
