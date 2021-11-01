//go:build !linux || !arm
// +build !linux !arm

package chamber

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
)

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

const stubTemperature = 25

func CreateThermometer(thermometerID string) (device.Thermometer, error) {
	return createStubThermometer(thermometerID)
}

func createStubThermometer(thermometerID string) (device.Thermometer, error) {
	return &stubThermometer{thermometerID: thermometerID}, nil
}

func CreateActuator(pin string) (device.Actuator, error) {
	return createStubActuator(pin)
}

func createStubActuator(pin string) (device.Actuator, error) {
	return &stubActuator{pin: pin}, nil
}

type stubThermometer struct {
	thermometerID string
}

func (t *stubThermometer) GetTemperature() (float64, error) {
	return stubTemperature, nil
}

type stubActuator struct {
	pin string
}

func (a *stubActuator) On() error { return nil }

func (a *stubActuator) Off() error { return nil }
