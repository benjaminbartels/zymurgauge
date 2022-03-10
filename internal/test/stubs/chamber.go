package stubs

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
)

const (
	stubTemperature = 25
	stubGravity     = 0.950
)

var (
	_ device.Thermometer              = (*Thermometer)(nil)
	_ device.Actuator                 = (*Actuator)(nil)
	_ device.ThermometerAndHydrometer = (*Tilt)(nil)
)

type Thermometer struct {
	ID string
}

func (t *Thermometer) GetID() string {
	return t.ID
}

func (t *Thermometer) GetTemperature() (float64, error) {
	return stubTemperature, nil
}

type Actuator struct {
	Pin string
}

func (a *Actuator) On() error { return nil }

func (a *Actuator) Off() error { return nil }

type Tilt struct {
	Color tilt.Color
}

func (t *Tilt) GetID() string {
	return string(t.Color)
}

func (t *Tilt) GetTemperature() (float64, error) {
	return stubTemperature, nil
}

func (t *Tilt) GetGravity() (float64, error) {
	return stubGravity, nil
}
