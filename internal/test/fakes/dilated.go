package fakes

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
)

var _ thermostat.Clock = (*DilatedClock)(nil)

type DilatedClock struct {
	multiplier float64
	startTime  time.Time
}

func NewDilatedClock(multiplier float64) thermostat.Clock {
	return &DilatedClock{
		multiplier: multiplier,
		startTime:  time.Now(),
	}
}

func (dc *DilatedClock) Now() time.Time {
	diff := float64(time.Since(dc.startTime)) / float64(time.Nanosecond)

	return dc.startTime.Add(time.Duration(dc.multiplier * diff))
}

func (dc *DilatedClock) After(d time.Duration) <-chan time.Time {
	d /= time.Duration(dc.multiplier)

	return time.After(d)
}

func (dc *DilatedClock) Since(t time.Time) time.Duration {
	return dc.startTime.Sub(t)
}

func (dc *DilatedClock) NewTimer(d time.Duration) *time.Timer {
	d /= time.Duration(dc.multiplier)

	return time.NewTimer(d)
}
