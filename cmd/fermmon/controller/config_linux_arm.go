package controller

import (
	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/ds18b20"
	"github.com/felixge/pidctrl"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

// ConfigureThermostat returns the results of ConfigurePiThermostat.  This methods only gets compiled when the operating
// system is linux and the architecture is arm.
func ConfigureThermostat(thermostat *internal.Thermostat, pid *pidctrl.PIDController,
	options ...internal.ThermostatOptionsFunc) error {
	return ConfigurePiThermostat(thermostat, pid, options...)
}

// Configure sets up a new chamber's RaspberryPI thermostat using the supplied ThermometerID, ChillerPin and HeaterPin
// values.
func ConfigurePiThermostat(thermostat *internal.Thermostat, pid *pidctrl.PIDController,
	options ...internal.ThermostatOptionsFunc) error {

	thermometer, err := ds18b20.NewThermometer(thermostat.ThermometerID)
	if err != nil {
		return err
	}

	adapter := raspi.NewAdaptor()

	var chiller *gpio.RelayDriver
	if thermostat.ChillerPin != "" {
		chiller = gpio.NewRelayDriver(adapter, thermostat.ChillerPin)
	}

	var heater *gpio.RelayDriver
	if thermostat.HeaterPin != "" {
		heater = gpio.NewRelayDriver(adapter, thermostat.HeaterPin)
	}

	// Setup the new Thermostat
	err = thermostat.Configure(pid, thermometer, chiller, heater, options...)
	if err != nil {
		return err
	}

	return nil

}
