package thermostat

import "time"

type Clock interface {
	Now() time.Time
}

// Real is a wrapper around Go's time package.
type RealClock struct{}

// New creates a new Clock.
func NewRealClock() Clock {
	return &RealClock{}
}

// Now returns the current local time.
func (*RealClock) Now() time.Time {
	return time.Now()
}
