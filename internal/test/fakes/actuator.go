package fakes

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/device"
)

var _ device.Actuator = (*Actuator)(nil)

type Actuator struct {
	OnFn      func() error
	OnCh      chan struct{}
	OffFn     func() error
	OffCh     chan time.Duration
	startTime time.Time
	OffError  error
}

func (a *Actuator) On() error {
	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}

	if a.OnCh != nil {
		a.OnCh <- struct{}{}
	}

	return nil
}

func (a *Actuator) Off() error {
	if a.OffError != nil {
		return a.OffError
	}

	if a.OffCh != nil {
		if a.startTime.IsZero() {
			a.OffCh <- 0
		} else {
			a.OffCh <- time.Since(a.startTime)
		}
	}

	return nil
}
