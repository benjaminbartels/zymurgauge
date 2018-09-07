// +build !linux !arm

package controller

import (
	"github.com/benjaminbartels/zymurgauge/internal"
)

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

// Configure sets up a new chamber's thermostat using the supplied ThermometerID, ChillerPin and HeaterPin values
func (c *ChamberCtl) Configure(chamber *internal.Chamber) error {
	return chamber.Thermostat.Configure(c.pid, &stubThermometer{}, &stubActuator{}, &stubActuator{},
		c.thermostatOptions...)
}

type stubThermometer struct{}

func (t *stubThermometer) Read() (*float64, error) {
	var f float64 = 22
	return &f, nil
}

type stubActuator struct{}

func (a *stubActuator) On() error { return nil }

func (a *stubActuator) Off() error { return nil }
