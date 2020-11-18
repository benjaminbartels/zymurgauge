package mocks

import "github.com/benjaminbartels/zymurgauge/internal/thermostat"

var _ thermostat.Actuator = (*Actuator)(nil)

type Actuator struct {
	OnFn  func() error
	OffFn func() error
}

func (a *Actuator) On() error {
	return a.OnFn()
}

func (a *Actuator) Off() error {
	return a.OffFn()
}
