package chamber

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
)

// TODO: move stubs to test folder

const (
	stubTemperature     = 25
	stubSpecificGravity = 0.950
)

var (
	_ device.Thermometer              = (*StubThermometer)(nil)
	_ device.Actuator                 = (*StubGPIOActuator)(nil)
	_ device.ThermometerAndHydrometer = (*StubTilt)(nil)
)

type StubThermometer struct {
	thermometerID string
}

func (t *StubThermometer) GetTemperature() (float64, error) {
	return stubTemperature, nil
}

type StubGPIOActuator struct {
	pin string
}

func (a *StubGPIOActuator) On() error { return nil }

func (a *StubGPIOActuator) Off() error { return nil }

type StubTilt struct {
	color tilt.Color
}

func (t *StubTilt) GetTemperature() (float64, error) {
	return stubTemperature, nil
}

func (t *StubTilt) GetSpecificGravity() (float64, error) {
	return stubSpecificGravity, nil
}
