//go:build !linux || !arm
// +build !linux !arm

package chamber

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

const stubTemperature = 25

func CreateThermometer(thermometerID string) (*StubThermometer, error) {
	return &StubThermometer{thermometerID: thermometerID}, nil
}

func CreateActuator(pin string) (*StubActuator, error) {
	return &StubActuator{pin: pin}, nil
}

type StubThermometer struct {
	thermometerID string
}

func (t *StubThermometer) GetTemperature() (float64, error) {
	return stubTemperature, nil
}

type StubActuator struct {
	pin string
}

func (a *StubActuator) On() error { return nil }

func (a *StubActuator) Off() error { return nil }
