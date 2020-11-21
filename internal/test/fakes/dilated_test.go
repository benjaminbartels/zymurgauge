package fakes_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/test/fakes"
)

func TestNow(t *testing.T) {
	tests := map[string]struct {
		multiplier float64
		round      time.Duration
	}{
		".01":   {multiplier: .01, round: 1 * time.Millisecond}, // 100ms = 1ms
		".1":    {multiplier: .1, round: 1 * time.Millisecond},  // 100ms = 10ms
		"10":    {multiplier: 10, round: 1 * time.Second},       // 100ms = 1s
		"100":   {multiplier: 100, round: 1 * time.Second},      // 100ms = 10s
		"600":   {multiplier: 600, round: 1 * time.Minute},      // 100ms = 1m
		"6000":  {multiplier: 6000, round: 1 * time.Minute},     // 100ms = 10m
		"36000": {multiplier: 36000, round: 1 * time.Hour},      // 100ms = 1h
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
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
	tests := map[string]struct {
		multiplier float64
		round      time.Duration
	}{
		".01":   {multiplier: .01, round: 1 * time.Millisecond},
		".1":    {multiplier: .1, round: 1 * time.Millisecond},
		"10":    {multiplier: 10, round: 1 * time.Second},
		"100":   {multiplier: 100, round: 1 * time.Second},
		"600":   {multiplier: 600, round: 1 * time.Minute},
		"6000":  {multiplier: 6000, round: 1 * time.Minute},
		"36000": {multiplier: 36000, round: 1 * time.Hour},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
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
	tests := map[string]struct {
		multiplier float64
		wait       time.Duration
	}{
		"10":    {multiplier: 10, wait: 1 * time.Second},    // 100ms = 1s
		"100":   {multiplier: 100, wait: 10 * time.Second},  // 100ms = 10s
		"600":   {multiplier: 600, wait: 1 * time.Minute},   // 100ms = 1m
		"6000":  {multiplier: 6000, wait: 10 * time.Minute}, // 100ms = 10m
		"36000": {multiplier: 36000, wait: 1 * time.Hour},   // 100ms = 1h
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			expected := 100 * time.Millisecond
			c := fakes.NewDilatedClock(tc.multiplier)
			timer := c.NewTimer(tc.wait)
			start := time.Now()

			<-timer.C // should be 100ms in real time

			since := time.Since(start)

			fmt.Println("since", since)

			if expected != since.Truncate(100*time.Millisecond) {
				t.Errorf("Unexpected diff. Want: '%s', Got: '%s'", expected, since.Truncate(100*time.Millisecond))
			}
		})
	}
}
