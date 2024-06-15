package tilt

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/pkg/errors"
)

var _ device.ThermometerAndHydrometer = (*Tilt)(nil)

type Tilt struct {
	color       Color
	temperature float64
	gravity     float64
	lastSeen    time.Time
}

var ErrIBeaconIsNil = errors.New("underlying IBeacon is nil")

func (t *Tilt) GetID() string {
	return string(t.color)
}

func (t *Tilt) GetTemperature() (float64, error) {
	return t.temperature, nil
}

func (t *Tilt) GetGravity() (float64, error) {
	gravityDenom := 1000.0

	return t.gravity / gravityDenom, nil
}
