package mocks

import "github.com/benjaminbartels/zymurgauge/internal/thermostat"

var _ thermostat.Thermometer = (*Thermometer)(nil)

type Thermometer struct {
	ReadFn func() (float64, error)
}

func (t *Thermometer) Read() (float64, error) {
	return t.ReadFn()
}
