package stubs

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
)

// TODO: move stubs to test folder

const (
	stubTemperature = 25
	stubGravity     = 0.950
)

var (
	_ device.Thermometer              = (*StubThermometer)(nil)
	_ device.Actuator                 = (*StubGPIOActuator)(nil)
	_ device.ThermometerAndHydrometer = (*StubTilt)(nil)
)

type StubThermometer struct {
	ThermometerID string
}

func (t *StubThermometer) GetTemperature() (float64, error) {
	return stubTemperature, nil
}

type StubGPIOActuator struct {
	Pin string
}

func (a *StubGPIOActuator) On() error { return nil }

func (a *StubGPIOActuator) Off() error { return nil }

type StubTilt struct {
	Color tilt.Color
}

func (t *StubTilt) GetTemperature() (float64, error) {
	return stubTemperature, nil
}

func (t *StubTilt) GetGravity() (float64, error) {
	return stubGravity, nil
}
