package main

import (
	"github.com/benjaminbartels/zymurgauge/internal/thermometer"
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func CreateThermostat(thermometer thermometer.Thermometer, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...thermostat.OptionsFunc) (*thermostat.Thermostat, error) {
	return CreatePiThermostat(thermometer, chillerPin, heaterPin, chillerKp, chillerKi, chillerKd, heaterKp, heaterKi,
		heaterKd, logger, options...)
}

func CreatePiThermostat(thermometer thermometer.Thermometer, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...thermostat.OptionsFunc) (*thermostat.Thermostat, error) {

	adapter := raspi.NewAdaptor()

	var chiller *gpio.RelayDriver
	if chillerPin != "" {
		chiller = gpio.NewRelayDriver(adapter, chillerPin)
	}

	var heater *gpio.RelayDriver
	if heaterPin != "" {
		heater = gpio.NewRelayDriver(adapter, heaterPin)
	}

	// Setup the new Thermostat
	return thermostat.NewThermostat(thermometer, chiller, heater, chillerKp, chillerKi, chillerKd, heaterKp,
		heaterKi, heaterKd, logger, options...), nil
}
