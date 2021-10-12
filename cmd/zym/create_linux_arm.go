package main

import (
	"github.com/benjaminbartels/zymurgauge/internal/platform/ds18b20"
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func CreateThermostat(thermometerID, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...thermostat.OptionsFunc) (*thermostat.Thermostat, error) {
	return CreatePiThermostat(thermometerID, chillerPin, heaterPin, chillerKp, chillerKi, chillerKd, heaterKp, heaterKi,
		heaterKd, logger, options...)
}

func CreatePiThermostat(thermometerID, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...thermostat.OptionsFunc) (*thermostat.Thermostat, error) {
	thermometer, err := ds18b20.NewThermometer(thermometerID)
	if err != nil {
		return nil, errors.Wrap(err, "could not creat enew thermometer")
	}

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
