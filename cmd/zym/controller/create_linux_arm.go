package controller

import (
	"github.com/benjaminbartels/zymurgauge/internal/device/pid"
	"github.com/benjaminbartels/zymurgauge/internal/device/raspberrypi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func CreateThermostat(thermometerID string, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...pid.OptionsFunc) (*pid.TemperatureController, error) {
	return createPiThermostat(thermometerID, chillerPin, heaterPin, chillerKp, chillerKi, chillerKd,
		heaterKp, heaterKi, heaterKd, logger, options...)
}

func createPiThermostat(thermometerID string, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...pid.OptionsFunc) (*pid.TemperatureController, error) {
	thermometer, err := raspberrypi.NewDs18b20(thermometerID)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new Ds18b20 thermometer")
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
	return pid.NewTemperatureController(thermometer, chiller, heater, chillerKp, chillerKi, chillerKd, heaterKp,
		heaterKi, heaterKd, logger, options...), nil
}
