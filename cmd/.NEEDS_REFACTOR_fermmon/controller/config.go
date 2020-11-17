// +build !linux !arm

package controller

import (
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/felixge/pidctrl"
)

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

// ConfigureThermostat returns the results of ConfigureStubThermostat.  This methods only gets compiled when the
// operating system is not linux and the architecture is not arm.
func ConfigureThermostat(thermostat *thermostat.Thermostat, pid *pidctrl.PIDController,
	options ...thermostat.ThermostatOptionsFunc) error {
	return ConfigureStubThermostat(thermostat, pid, options...)
}
