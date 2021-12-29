package configurator

import "github.com/benjaminbartels/zymurgauge/internal/device/tilt"

type ConfiguratorIface interface {
	CreateDs18b20(thermometerID string) (*StubThermometer, error)
	CreateTilt(color tilt.Color) (*StubTilt, error)
	CreateGPIOActuator(pin string) (*StubGPIOActuator, error)
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

func (c StubConfigurator) CreateDs18b20(thermometerID string) (*StubThermometer, error) {
	return &StubThermometer{thermometerID: thermometerID}, nil
}

func (c StubConfigurator) CreateTilt(color tilt.Color) (*StubTilt, error) {
	return &StubTilt{color: color}, nil
}

func (c StubConfigurator) CreateGPIOActuator(pin string) (*StubGPIOActuator, error) {
	return &StubGPIOActuator{pin: pin}, nil
}
