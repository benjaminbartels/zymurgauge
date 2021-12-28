package clock

import "time"

var _ Clock = (*RealClockPPP)(nil)

type Clock interface {
	Now() time.Time
	Since(t time.Time) time.Duration
	NewTimer(d time.Duration) *time.Timer
}

// Real is a wrapper around Go's time package.
type RealClockPPP struct{}

// New creates a new Clock.
func NewRealClock() *RealClockPPP {
	return &RealClockPPP{}
}

// Now returns the current local time.
func (*RealClockPPP) Now() time.Time {
	return time.Now()
}

// Since returns the time elapsed since t.
func (*RealClockPPP) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// NewTimer creates a new Timer that will send
// the current time on its channel after at least duration d.s.
func (*RealClockPPP) NewTimer(d time.Duration) *time.Timer {
	return time.NewTimer(d)
}
