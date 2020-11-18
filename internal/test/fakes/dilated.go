package fakes

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
)

var _ thermostat.Clock = (*DilatedClock)(nil)

type DilatedClock struct {
	speed     float64
	startTime time.Time
}

func NewDilatedClock(speed float64) thermostat.Clock {
	return &DilatedClock{
		speed:     speed,
		startTime: time.Now(),
	}
}

func (dc *DilatedClock) Now() time.Time {
	diff := float64(time.Since(dc.startTime)) / float64(time.Nanosecond)

	return dc.startTime.Add(time.Duration(dc.speed * diff))
}
