package tilt

import (
	"math"

	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/platform/bluetooth"
	"github.com/pkg/errors"
)

var _ device.ThermometerAndHydrometer = (*Tilt)(nil)

type Tilt struct {
	ibeacon *bluetooth.IBeacon
	color   Color
}

var ErrIBeaconIsNil = errors.New("underlying IBeacon is nil")

func (t *Tilt) GetID() string {
	return string(t.color)
}

func (t *Tilt) GetTemperature() (float64, error) {
	if t.ibeacon == nil {
		return 0, ErrIBeaconIsNil
	}

	return math.Round(float64(t.ibeacon.GetMajor()-32)/1.8*100) / 100, nil
}

func (t *Tilt) GetGravity() (float64, error) {
	if t.ibeacon == nil {
		return 0, ErrIBeaconIsNil
	}

	return float64(t.ibeacon.GetMinor()) / 1000, nil
}
