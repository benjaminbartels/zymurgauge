package fakes

import (
	"testing"
	"time"
)

func TestDilation(t *testing.T) {
	tests := map[string]struct {
		multiplier float64
		round      time.Duration
		diff       time.Duration
	}{
		".01":   {multiplier: .01, round: 1 * time.Millisecond, diff: 1 * time.Millisecond}, // 100ms = 1ms
		".1":    {multiplier: .1, round: 1 * time.Millisecond, diff: 10 * time.Millisecond}, // 100ms = 10ms
		"10":    {multiplier: 10, round: 1 * time.Second, diff: 1 * time.Second},            // 100ms = 1s
		"100":   {multiplier: 100, round: 1 * time.Second, diff: 10 * time.Second},          // 100ms = 10s
		"600":   {multiplier: 600, round: 1 * time.Minute, diff: 1 * time.Minute},           // 100ms = 1m
		"6000":  {multiplier: 6000, round: 1 * time.Minute, diff: 10 * time.Minute},         // 100ms = 10m
		"36000": {multiplier: 36000, round: 1 * time.Hour, diff: 1 * time.Hour},             // 100ms = 1h
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			start := time.Now()
			c := NewDilatedClock(tc.multiplier)

			<-time.After(100 * time.Millisecond)

			dilatedEnd := c.Now()
			dilatedDiff := dilatedEnd.Sub(start)

			if tc.diff != dilatedDiff.Truncate(tc.round) {
				t.Errorf("Unexpected dilated diff. Wanted '%s', Got: '%s'", tc.diff, dilatedDiff.Truncate(tc.round))
			}
		})
	}
}
