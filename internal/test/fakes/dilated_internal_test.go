package fakes

import (
	"testing"
	"time"
)

func TestDilation(t *testing.T) {
	tests := map[string]struct {
		speed float64
		round time.Duration
		diff  time.Duration
	}{
		".01":   {speed: .01, round: 1 * time.Millisecond, diff: 1 * time.Millisecond}, // 100ms = 1ms
		".1":    {speed: .1, round: 1 * time.Millisecond, diff: 10 * time.Millisecond}, // 100ms = 10ms
		"10":    {speed: 10, round: 1 * time.Second, diff: 1 * time.Second},            // 100ms = 1s
		"100":   {speed: 100, round: 1 * time.Second, diff: 10 * time.Second},          // 100ms = 10s
		"600":   {speed: 600, round: 1 * time.Minute, diff: 1 * time.Minute},           // 100ms = 1m
		"6000":  {speed: 6000, round: 1 * time.Minute, diff: 10 * time.Minute},         // 100ms = 10m
		"36000": {speed: 36000, round: 1 * time.Hour, diff: 1 * time.Hour},             // 100ms = 1h
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			start := time.Now()
			c := NewDilatedClock(tc.speed)

			<-time.After(100 * time.Millisecond)

			dilatedEnd := c.Now()
			dilatedDiff := dilatedEnd.Sub(start)

			if tc.diff != dilatedDiff.Truncate(tc.round) {
				t.Errorf("Unexpected dilated diff. Wanted '%s', Got: '%s'", tc.diff, dilatedDiff.Truncate(tc.round))
			}
		})
	}
}
