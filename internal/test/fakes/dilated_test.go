package fakes_test

import (
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/test/fakes"
)

type dilationTestParams struct {
	name       string
	multiplier float64
	round      time.Duration
}

func get01to36000Tests() []dilationTestParams {
	return []dilationTestParams{
		{name: ".01", multiplier: .01, round: 1 * time.Millisecond}, // 100ms = 1ms
		{name: ".1", multiplier: .1, round: 1 * time.Millisecond},   // 100ms = 10ms
		{name: "10", multiplier: 10, round: 1 * time.Second},        // 100ms = 1s
		{name: "100", multiplier: 100, round: 1 * time.Second},      // 100ms = 10s
		{name: "600", multiplier: 600, round: 1 * time.Minute},      // 100ms = 1m
		{name: "6000", multiplier: 6000, round: 1 * time.Minute},    // 100ms = 10m
		{name: "36000", multiplier: 36000, round: 1 * time.Hour},    // 100ms = 1h
	}
}

func TestNow(t *testing.T) {
	t.Parallel()

	tests := get01to36000Tests()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := fakes.NewDilatedClock(tc.multiplier)
			start := time.Now()

			<-time.After(100 * time.Millisecond)

			expected := c.Since(start)
			diff := c.Now().Sub(start)

			if expected.Truncate(tc.round) != diff.Truncate(tc.round) {
				t.Errorf("Unexpected diff. Want: '%s', Got: '%s'", expected.Truncate(tc.round), diff.Truncate(tc.round))
			}
		})
	}
}

func TestSince(t *testing.T) {
	t.Parallel()

	tests := get01to36000Tests()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := fakes.NewDilatedClock(tc.multiplier)
			start := c.Now()

			<-time.After(100 * time.Millisecond)

			since := c.Since(start)
			expected := c.Now().Sub(start)

			if expected.Truncate(tc.round) != since.Truncate(tc.round) {
				t.Errorf("Unexpected diff. Want: '%s', Got: '%s'", expected.Truncate(tc.round), since.Truncate(tc.round))
			}
		})
	}
}

func TestNewTimer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		multiplier float64
		wait       time.Duration
	}{
		{name: "10", multiplier: 10, wait: 1 * time.Second},      // 100ms = 1s
		{name: "100", multiplier: 100, wait: 10 * time.Second},   // 100ms = 10s
		{name: "600", multiplier: 600, wait: 1 * time.Minute},    // 100ms = 1m
		{name: "6000", multiplier: 6000, wait: 10 * time.Minute}, // 100ms = 10m
		{name: "36000", multiplier: 36000, wait: 1 * time.Hour},  // 100ms = 1h
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			expected := 100 * time.Millisecond
			c := fakes.NewDilatedClock(tc.multiplier)
			timer := c.NewTimer(tc.wait)
			start := time.Now()

			<-timer.C // should be 100ms in real time

			since := time.Since(start)

			if expected != since.Truncate(100*time.Millisecond) {
				t.Errorf("Unexpected diff. Want: '%s', Got: '%s'", expected, since.Truncate(100*time.Millisecond))
			}
		})
	}
}
