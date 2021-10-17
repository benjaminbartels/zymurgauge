package pid

import "time"

var _ Clock = (*RealClock)(nil)

type Clock interface {
	Now() time.Time
	Since(t time.Time) time.Duration
	NewTimer(d time.Duration) *time.Timer
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

// Since returns the time elapsed since t.
func (*RealClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// NewTimer creates a new Timer that will send
// the current time on its channel after at least duration d.s.
func (*RealClock) NewTimer(d time.Duration) *time.Timer {
	return time.NewTimer(d)
}
