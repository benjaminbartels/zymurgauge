package main

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/client"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/felixge/pidctrl"
)

// Chamber represents an insulated box (fridge) with internal heating/cooling elements that reacts to changes in
// monitored temperatures, by correcting small deviations from your desired fermentation temperature.
type Chamber struct {
	chamberDefinition   *internal.Chamber
	pid                 *pidctrl.PIDController
	client              *client.Client
	logger              log.Logger
	thermostatOptions   []func(*Thermostat) error
	thermostat          *Thermostat
	currentFermentation *internal.Fermentation
}

// NewChamber creates a new Chamber with the given parameters
func NewChamber(chamberDefinition *internal.Chamber, pid *pidctrl.PIDController, client *client.Client,
	logger log.Logger, thermostatOptions ...func(*Thermostat) error) (*Chamber, error) {

	var err error

	c := &Chamber{
		chamberDefinition: chamberDefinition,
		pid:               pid,
		client:            client,
		logger:            logger,
		thermostatOptions: thermostatOptions,
	}

	c.thermostat, err = NewThermostat(chamberDefinition.Thermostat, pid, logger, thermostatOptions...)
	if err != nil {
		return nil, err
	}

	if chamberDefinition.CurrentFermentationID != 0 {
		c.currentFermentation, err = client.FermentationResource.Get(chamberDefinition.CurrentFermentationID)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Listen subscribes to receive Chamber updates from the server
func (c *Chamber) Listen() {

	// Subscribe to chamber updates
	ch := make(chan internal.Chamber)
	err := c.client.ChamberResource.Subscribe(c.chamberDefinition.MacAddress, ch)
	if err != nil {
		c.logger.Fatal(err) // ToDo Retry on error?
	}

	for { // To Do: cancel?
		c.logger.Println("Waiting for ChamberService updates")
		select {
		case cham := <-ch:
			if err = c.processUpdate(&cham); err != nil {
				c.logger.Println(err.Error()) // ToDo: Handle
			}
			// case sig := <-signals:
			// 	fmt.Println("received signal", sig)
		}
	}
}

func (c *Chamber) processUpdate(chamberDefinition *internal.Chamber) error {
	c.logger.Println("Processing Chamber:", chamberDefinition.Name, chamberDefinition.MacAddress)

	var err error

	if c.chamberDefinition.Thermostat != chamberDefinition.Thermostat {
		c.thermostat.Off()

		// Create a new Thermostat
		c.thermostat, err = NewThermostat(chamberDefinition.Thermostat, c.pid, c.logger, c.thermostatOptions...)
		if err != nil {
			return err
		}

		c.thermostat.Subscribe(chamberDefinition.MacAddress, c.handleStatusUpdate)
	}

	if c.currentFermentation.ID != chamberDefinition.CurrentFermentationID {
		c.thermostat.Off()

		// If CurrentFermentationID != 0 then get the fermentation
		if chamberDefinition.CurrentFermentationID != 0 {
			c.currentFermentation, err = c.client.FermentationResource.Get(chamberDefinition.CurrentFermentationID)
			if err != nil {
				return err
			}
		} else {
			c.logger.Println("No current Fermentation")
			c.currentFermentation = nil
		}
	}

	// If there is a fermentation and thermostat is not nil, change setting and turn thermostat On if it is Off
	if c.currentFermentation != nil && c.thermostat != nil {
		c.thermostat.Set(c.currentFermentation.Beer.Schedule[0].TargetTemp)
		if c.thermostat.GetStatus().State == OFF {
			c.thermostat.On()
		}
	}

	// set the definition
	c.chamberDefinition = chamberDefinition

	return nil
}

func (c *Chamber) handleStatusUpdate(status ThermostatStatus) {
	c.logger.Printf("State: %v, Error: %s\n", status.State, status.Error)

	change := &internal.TemperatureChange{
		FermentationID: c.currentFermentation.ID,
		InsertTime:     time.Now(),
		Chamber:        c.chamberDefinition.Name,
		Beer:           c.currentFermentation.Beer.Name,
		Thermometer:    c.thermostat.thermometer.ID,
		Temperature:    *status.CurrentTemperature,
	}
	if err := c.client.FermentationResource.SaveTemperatureChange(change); err != nil {
		c.logger.Println(err)
	}
}
