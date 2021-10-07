package mocks

import "github.com/benjaminbartels/zymurgauge/internal/thermometer"

var _ thermometer.Thermometer = (*Thermometer)(nil)

type Thermometer struct {
	ReadFn func() (float64, error)
}

func (t *Thermometer) Read() (float64, error) {
	return t.ReadFn()
}
