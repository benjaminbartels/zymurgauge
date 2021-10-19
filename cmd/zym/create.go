//go:build !linux || !arm
// +build !linux !arm

package main

import (
	"github.com/benjaminbartels/zymurgauge/internal/pid"
	"github.com/sirupsen/logrus"
)

// The program is only meant to run on linux on arm. This file only exists to prevent compilation issues on non
// linux/arm systems.

const stubTemperature = 22

func CreateTemperatureController(thermometerAddress uint64, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...pid.OptionsFunc) (*pid.TemperatureController, error) {
	return CreateStubTemperatureController(thermometerAddress, chillerPin, heaterPin, chillerKp, chillerKi, chillerKd,
		heaterKp, heaterKi, heaterKd, logger, options...)
}

func CreateStubTemperatureController(thermometerAddress uint64, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...pid.OptionsFunc) (*pid.TemperatureController, error) {
	return pid.NewTemperatureController(&stubThermometer{thermometerAddress}, &stubActuator{chillerPin},
		&stubActuator{heaterPin}, chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd, logger,
		options...), nil
}

type stubThermometer struct {
	thermometerAddress uint64
}

func (t *stubThermometer) GetTemperature() (float64, error) {
	return stubTemperature, nil
}

type stubActuator struct {
	pin string
}

func (a *stubActuator) On() error { return nil }

func (a *stubActuator) Off() error { return nil }
