// +build !linux !arm

package controller

import (
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/sirupsen/logrus"
)

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

func CreateThermostat(thermometerID, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...thermostat.OptionsFunc) (*thermostat.Thermostat, error) {
	return CreateStubThermostat(thermometerID, chillerPin, heaterPin, chillerKp, chillerKi, chillerKd, heaterKp,
		heaterKi, heaterKd, logger, options...)
}
