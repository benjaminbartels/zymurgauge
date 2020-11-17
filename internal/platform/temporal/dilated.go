package temporal

import (
	"time"
)

type DilatedClock struct {
	speed     float64
	startTime time.Time
}

func NewDilatedClock(speed float64) Clock {
	return &DilatedClock{
		speed:     speed,
		startTime: time.Now(),
	}
}

func (dc *DilatedClock) Now() time.Time {
	diff := float64(time.Since(dc.startTime)) / float64(time.Nanosecond)
	return dc.startTime.Add(time.Duration(dc.speed * diff))
}

func (dc *DilatedClock) After(d time.Duration) <-chan time.Time {
	d = d / time.Duration(dc.speed)
	return time.After(d)
}
