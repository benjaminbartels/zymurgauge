package main

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/client"
	"github.com/benjaminbartels/zymurgauge/internal/gpio"
	"github.com/benjaminbartels/zymurgauge/internal/platform/ds18b20"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/felixge/pidctrl"
)

// ChamberCtl controls a Chamber
type chamberCtl struct {
	chamber             *internal.Chamber
	pid                 *pidctrl.PIDController
	client              *client.Client
	logger              log.Logger
	thermostatOptions   []func(*internal.Thermostat) error
	currentFermentation *internal.Fermentation
	done                chan bool
}

// NewChamberCtl creates a new ChamberCtl with the given parameters
func newChamberCtl(chamber *internal.Chamber, pid *pidctrl.PIDController, client *client.Client,
	logger log.Logger, thermostatOptions ...func(*internal.Thermostat) error) (*chamberCtl, error) {

	var err error

	c := &chamberCtl{
		chamber:           chamber,
		pid:               pid,
		client:            client,
		logger:            logger,
		thermostatOptions: thermostatOptions,
		done:              make(chan bool, 1),
	}

	thermometer, err := ds18b20.NewThermometer(c.chamber.Thermostat.ThermometerID)
	if err != nil {
		return nil, err
	}

	chiller, err := gpio.NewActuator(c.chamber.Thermostat.ChillerPin)
	if err != nil {
		return nil, err
	}

	heater, err := gpio.NewActuator(c.chamber.Thermostat.HeaterPin)
	if err != nil {
		return nil, err
	}

	err = c.chamber.Thermostat.Setup(pid, thermometer, chiller, heater, thermostatOptions...)
	if err != nil {
		return nil, err
	}

	if c.chamber.CurrentFermentationID != 0 {
		c.currentFermentation, err = client.FermentationResource.Get(c.chamber.CurrentFermentationID)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Listen subscribes to receive Chamber updates from the server
func (c *chamberCtl) Listen() {

	// Subscribe to chamber updates
	ch := make(chan internal.Chamber)
	err := c.client.ChamberResource.Subscribe(c.chamber.MacAddress, ch)
	if err != nil {
		c.logger.Fatal(err) // ToDo Retry on error?
	}

	for {
		c.logger.Println("Waiting for ChamberService updates")
		select {
		case cham := <-ch:
			if err = c.processUpdate(&cham); err != nil {
				c.logger.Println(err.Error()) // ToDo: Handle
			}
		case <-c.done:
			c.client.ChamberResource.Unsubscribe(c.chamber.MacAddress)
		}
	}
}

// Close unsubscribes from receiving Chamber updates from the server
func (c *chamberCtl) Close() {
	c.done <- true
}

func (c *chamberCtl) processUpdate(chamberUpdate *internal.Chamber) error {
	c.logger.Println("Processing Chamber:", chamberUpdate.Name, chamberUpdate.MacAddress)

	var err error

	if c.chamber.Thermostat.ChillerPin != chamberUpdate.Thermostat.ChillerPin ||
		c.chamber.Thermostat.HeaterPin != chamberUpdate.Thermostat.HeaterPin ||
		c.chamber.Thermostat.ThermometerID != chamberUpdate.Thermostat.ThermometerID {

		c.chamber.Thermostat.Off()

		thermometer, err := ds18b20.NewThermometer(c.chamber.Thermostat.ThermometerID)
		if err != nil {
			return err
		}

		chiller, err := gpio.NewActuator(c.chamber.Thermostat.ChillerPin)
		if err != nil {
			return err
		}

		heater, err := gpio.NewActuator(c.chamber.Thermostat.HeaterPin)
		if err != nil {
			return err
		}

		// Setup the new Thermostat
		err = chamberUpdate.Thermostat.Setup(c.pid, thermometer, chiller, heater, c.thermostatOptions...)
		if err != nil {
			return err
		}

		c.chamber.Thermostat.Subscribe(chamberUpdate.MacAddress, c.handleStatusUpdate)
	}

	if c.currentFermentation.ID != chamberUpdate.CurrentFermentationID {
		c.chamber.Thermostat.Off()

		// If CurrentFermentationID != 0 then get the fermentation
		if chamberUpdate.CurrentFermentationID != 0 {
			c.currentFermentation, err = c.client.FermentationResource.Get(chamberUpdate.CurrentFermentationID)
			if err != nil {
				return err
			}
		} else {
			c.logger.Println("No current Fermentation")
			c.currentFermentation = nil
		}
	}

	// If there is a fermentation and thermostat is not nil, change setting and turn thermostat On if it is Off
	if c.currentFermentation != nil { // To Do: remove && c.chamber.Thermostat != nil {
		c.chamber.Thermostat.Set(c.currentFermentation.Beer.Schedule[0].TargetTemp)
		if c.chamber.Thermostat.GetStatus().State == internal.OFF {
			c.chamber.Thermostat.On()
		}
	}

	// set the definition
	c.chamber = chamberUpdate

	return nil
}

func (c *chamberCtl) handleStatusUpdate(status internal.ThermostatStatus) {
	c.logger.Printf("State: %v, Error: %s\n", status.State, status.Error)

	change := &internal.TemperatureChange{
		FermentationID: c.currentFermentation.ID,
		InsertTime:     time.Now(),
		Chamber:        c.chamber.Name,
		Beer:           c.currentFermentation.Beer.Name,
		Thermometer:    c.chamber.Thermostat.ThermometerID,
		Temperature:    *status.CurrentTemperature,
	}
	if err := c.client.FermentationResource.SaveTemperatureChange(change); err != nil {
		c.logger.Println(err)
	}
}
