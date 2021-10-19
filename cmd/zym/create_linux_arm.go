package main

import (
	"github.com/benjaminbartels/zymurgauge/internal/pid"
	"github.com/benjaminbartels/zymurgauge/internal/thermometer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"periph.io/x/periph/conn/onewire"
	"periph.io/x/periph/experimental/host/netlink"
)

const resolutionBits = 10

func CreateThermostat(thermometerAddress uint64, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...pid.OptionsFunc) (*pid.TemperatureController, error) {
	return CreatePiThermostat(thermometerAddress, chillerPin, heaterPin, chillerKp, chillerKi, chillerKd,
		heaterKp, heaterKi, heaterKd, logger, options...)
}

func CreatePiThermostat(thermometerAddress uint64, chillerPin, heaterPin string,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...pid.OptionsFunc) (*pid.TemperatureController, error) {

	bus, err := netlink.New(001)
	if err != nil {
		return nil, errors.Wrap(err, "could not open 1-wire bus")
	}

	thermometer, err := thermometer.NewDs18b20(bus, onewire.Address(thermometerAddress), resolutionBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new TemperatureController")
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
