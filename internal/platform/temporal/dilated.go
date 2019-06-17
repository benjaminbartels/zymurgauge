package temporal

import (
	"time"
)

type DilatedClock struct {
	multiplyer float64
	startTime  time.Time
}

func NewDilatedClock(multiplyer float64) Clock {
	return &DilatedClock{
		multiplyer: multiplyer,
		startTime:  time.Now(),
	}
}

func (dc *DilatedClock) Now() time.Time {
	diff := float64(time.Since(dc.startTime)) / float64(time.Nanosecond)
	return dc.startTime.Add(time.Duration(dc.multiplyer * diff))
}

func (dc *DilatedClock) After(d time.Duration) <-chan time.Time {
	d = d / time.Duration(dc.multiplyer)
	return time.After(d)
}
