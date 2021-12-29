package chamber

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
)

type Configurator interface {
	CreateDs18b20(thermometerID string) (device.Thermometer, error)
	CreateTilt(color tilt.Color) (device.ThermometerAndHydrometer, error)
	CreateGPIOActuator(pin string) (device.Actuator, error)
}

const (
	stubTemperature     = 25
	stubSpecificGravity = 0.950
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

type StubConfigurator struct{}

func (c StubConfigurator) CreateDs18b20(thermometerID string) (device.Thermometer, error) {
	return &StubThermometer{thermometerID: thermometerID}, nil
}

func (c StubConfigurator) CreateTilt(color tilt.Color) (device.ThermometerAndHydrometer, error) {
	return &StubTilt{color: color}, nil
}

func (c StubConfigurator) CreateGPIOActuator(pin string) (device.Actuator, error) {
	return &StubGPIOActuator{pin: pin}, nil
}
