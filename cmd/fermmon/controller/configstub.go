package controller

import (
	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/felixge/pidctrl"
)

// ConfigureStubThermostat configures a thermostat using stubs for the Thermometer and Actuators.
func ConfigureStubThermostat(thermostat *internal.Thermostat, pid *pidctrl.PIDController,
	options ...internal.ThermostatOptionsFunc) error {
	return thermostat.Configure(pid, &stubThermometer{}, &stubActuator{}, &stubActuator{}, options...)
}

type stubThermometer struct{}

func (t *stubThermometer) Read() (*float64, error) {
	var f float64 = 22
	return &f, nil
}

type stubActuator struct{}

func (a *stubActuator) On() error { return nil }

func (a *stubActuator) Off() error { return nil }
