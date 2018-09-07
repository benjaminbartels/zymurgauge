package controller

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/web"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/client"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/felixge/pidctrl"
	"github.com/pkg/errors"
)

const interval = 10 * time.Second

// ChamberCtl controls a Chamber
type ChamberCtl struct {
	mac               string
	chamber           *internal.Chamber
	pid               *pidctrl.PIDController
	client            *client.Client
	logger            log.Logger
	thermostatOptions []func(*internal.Thermostat) error
	fermentation      *internal.Fermentation
	done              chan bool
}

// NewChamberCtl creates a new ChamberCtl with the given parameters
func NewChamberCtl(mac string, pid *pidctrl.PIDController, client *client.Client,
	logger log.Logger, thermostatOptions ...func(*internal.Thermostat) error) *ChamberCtl {

	return &ChamberCtl{
		mac:               mac,
		pid:               pid,
		client:            client,
		logger:            logger,
		thermostatOptions: thermostatOptions,
		done:              make(chan bool, 1),
	}

}

// Start begin the process of polling the service for updates to the chamber
func (c *ChamberCtl) Start() {

	c.poll()

	for {
		select {
		case <-c.done:
			return
		case <-time.After(interval):
			c.poll()
		}
	}
}

func (c *ChamberCtl) poll() {
	var chamber *internal.Chamber
	chamber, err := c.client.ChamberResource.Get(c.mac)
	if err != nil {
		if web.ErrNotFound == errors.Cause(err) {
			c.logger.Println("Chamber does not exist. Creating new chamber")
			chamber := &internal.Chamber{
				Name:       "Chamber " + c.mac,
				MacAddress: c.mac,
			}
			err := c.client.ChamberResource.Save(chamber)
			if err != nil {
				c.logger.Fatalln(err.Error()) // ToDo: Handle
			}
		}
	} else {
		if err = c.processUpdate(chamber); err != nil {
			c.logger.Println(err.Error()) // ToDo: Handle
		}
	}
}

// Stop ends the polling process
func (c *ChamberCtl) Stop() {
	c.done <- true
}

// processUpdate evaluates the inbound chamber to determine if any changes have occurred
func (c *ChamberCtl) processUpdate(chamber *internal.Chamber) error {

	configChanged := false

	if c.chamber == nil {
		configChanged = true
		c.chamber = chamber
	} else {

		if c.chamber.Thermostat.ChillerPin != chamber.Thermostat.ChillerPin {
			c.logger.Println("Chiller Pin changed from %s to %s",
				c.chamber.Thermostat.ChillerPin, chamber.Thermostat.ChillerPin)

			configChanged = true
		}

		if c.chamber.Thermostat.HeaterPin != chamber.Thermostat.HeaterPin {
			c.logger.Println("Heater Pin changed from %s to %s",
				c.chamber.Thermostat.HeaterPin, chamber.Thermostat.HeaterPin)

			configChanged = true
		}

		if c.chamber.Thermostat.ChillerPin != chamber.Thermostat.ChillerPin {
			c.logger.Println("Thermometer ID changed from %s to %s",
				c.chamber.Thermostat.ChillerPin, chamber.Thermostat.ChillerPin)

			configChanged = true
		}
	}

	if configChanged {

		c.chamber.Thermostat.Off()

		c.chamber = chamber

		if err := c.Configure(c.chamber); err != nil {
			return err
		}

		c.chamber.Thermostat.Subscribe(c.chamber.MacAddress, c.handleStatusUpdate)
	}

	var err error
	if c.fermentation == nil {
		if c.chamber.CurrentFermentationID != 0 {
			c.logger.Println("Fermentation changed from none to %d", c.chamber.CurrentFermentationID)
			c.fermentation, err = c.getFermentation(c.chamber.CurrentFermentationID)
			if err != nil {
				return err
			}
		}
	} else {
		if c.chamber.CurrentFermentationID != 0 {
			c.logger.Println("Fermentation changed from %d to %d", c.chamber.CurrentFermentationID, c.fermentation.ID)
			c.chamber.Thermostat.Off()
			c.fermentation, err = c.getFermentation(c.chamber.CurrentFermentationID)
			if err != nil {
				return err
			}
		} else {
			c.logger.Println("Fermentation changed from %d to none", c.chamber.CurrentFermentationID)
			c.chamber.Thermostat.Off()
		}
	}

	if c.fermentation != nil {
		c.logger.Printf("Current Fermentation has be set to %d\n", c.fermentation.ID)
		c.chamber.Thermostat.Set(c.fermentation.Beer.Schedule[0].TargetTemp)
		if c.chamber.Thermostat.GetStatus().State == internal.OFF {
			c.chamber.Thermostat.On()
		}
	}

	return nil
}

func (c *ChamberCtl) getFermentation(id uint64) (*internal.Fermentation, error) {

	fermentation, err := c.client.FermentationResource.Get(id)
	if err != nil {
		return nil, err
	}

	return fermentation, nil
}

func (c *ChamberCtl) handleStatusUpdate(status internal.ThermostatStatus) {

	var temp float64
	var errMsg string

	if status.CurrentTemperature != nil {
		temp = *status.CurrentTemperature
	}

	if status.Error != nil {
		errMsg = status.Error.Error()
	}

	c.logger.Printf("State: %v, Temperature: %f, Error: %s\n", status.State, temp, errMsg)

	change := &internal.TemperatureChange{
		InsertTime:  time.Now(),
		Chamber:     c.chamber.Name,
		Thermometer: c.chamber.Thermostat.ThermometerID,
	}

	if c.fermentation != nil {
		change.FermentationID = c.fermentation.ID
		change.Beer = c.fermentation.Beer.Name
	}

	if status.CurrentTemperature != nil {
		change.Temperature = *status.CurrentTemperature
	}

	if err := c.client.FermentationResource.SaveTemperatureChange(change); err != nil {
		c.logger.Println(err)
	}
}
