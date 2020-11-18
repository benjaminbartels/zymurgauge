package controller

import (
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/sirupsen/logrus"
)

const temp = 22

func CreateStubThermostat(thermometerID, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...thermostat.OptionsFunc) (*thermostat.Thermostat, error) {
	return thermostat.NewThermostat(&stubThermometer{}, &stubActuator{}, &stubActuator{},
		chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd, logger,
		options...), nil
}

type stubThermometer struct{}

func (t *stubThermometer) Read() (float64, error) {
	return temp, nil
}

type stubActuator struct{}

func (a *stubActuator) On() error { return nil }

func (a *stubActuator) Off() error { return nil }
