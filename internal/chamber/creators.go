//go:build !linux || !arm
// +build !linux !arm

package chamber

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

// TODO: Force stubs from unit tests?

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
	color string
}

func (t *StubTilt) GetTemperature() (float64, error) {
	return stubTemperature, nil
}

func (t *StubTilt) GetSpecificGravity() (float64, error) {
	return stubSpecificGravity, nil
}

func CreateDs18b20(thermometerID string) (*StubThermometer, error) {
	return &StubThermometer{thermometerID: thermometerID}, nil
}

func CreateTilt(color string) (*StubTilt, error) {
	return &StubTilt{color: color}, nil
}

func CreateGPIOActuator(pin string) (*StubGPIOActuator, error) {
	return &StubGPIOActuator{pin: pin}, nil
}
