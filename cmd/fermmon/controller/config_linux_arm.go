package controller

import (
	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/ds18b20"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

// Configure sets up a new chamber's thermostat using the supplied ThermometerID, ChillerPin and HeaterPin values
func (c *ChamberCtl) Configure(chamber *internal.Chamber) error {

	thermometer, err := ds18b20.NewThermometer(chamber.Thermostat.ThermometerID)
	if err != nil {
		return err
	}

	adapter := raspi.NewAdaptor()

	var chiller *gpio.RelayDriver
	if c.chamber.Thermostat.ChillerPin != "" {
		chiller = gpio.NewRelayDriver(adapter, chamber.Thermostat.ChillerPin)
	}

	var heater *gpio.RelayDriver
	if c.chamber.Thermostat.HeaterPin != "" {
		heater = gpio.NewRelayDriver(adapter, chamber.Thermostat.HeaterPin)
	}

	// Setup the new Thermostat
	err = chamber.Thermostat.Configure(c.pid, thermometer, chiller, heater, c.thermostatOptions...)
	if err != nil {
		return err
	}

	return nil

}
