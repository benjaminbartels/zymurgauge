package main

import (
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/sirupsen/logrus"
)

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

const temp = 22

func CreateThermostat(thermometerID, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...thermostat.OptionsFunc) (*thermostat.Thermostat, error) {
	return CreateStubThermostat(thermometerID, chillerPin, heaterPin, chillerKp, chillerKi, chillerKd, heaterKp,
		heaterKi, heaterKd, logger, options...)
}

func CreateStubThermostat(thermometerID, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...thermostat.OptionsFunc) (*thermostat.Thermostat, error) {
	return thermostat.NewThermostat(&stubThermometer{thermometerID}, &stubActuator{chillerPin}, &stubActuator{heaterPin},
		chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd, logger,
		options...), nil
}

type stubThermometer struct {
	thermometerID string
}

func (t *stubThermometer) Read() (float64, error) {
	return temp, nil
}

type stubActuator struct {
	pin string
}

func (a *stubActuator) On() error { return nil }

func (a *stubActuator) Off() error { return nil }
