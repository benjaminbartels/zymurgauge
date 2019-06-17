package clock

import "time"

// Real is a wrapper around Go's time package.
type RealClock struct{}

// New creates a new Clock
func New() Clock {
	return &RealClock{}
}

// Now returns the current local time.
func (*RealClock) Now() time.Time {
	return time.Now()
}

// After waits for the duration to elapse and then sends the current time
// on the returned channel.
func (*RealClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}
