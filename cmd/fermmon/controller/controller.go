package controller

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/client"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/felixge/pidctrl"
	"github.com/pkg/errors"
)

// Controller controls a Chamber. ConfigFunc is the function used to configure the controller's Thermostat.
type Controller struct {
	ConfigFunc           func(*internal.Thermostat, *pidctrl.PIDController, ...internal.ThermostatOptionsFunc) error
	Chamber              *internal.Chamber
	Fermentation         *internal.Fermentation
	id                   string
	mac                  string
	pid                  *pidctrl.PIDController
	chamberProvider      client.ChamberProvider
	fermentationProvider client.FermentationProvider
	logger               log.Logger
	thermostatOptions    []internal.ThermostatOptionsFunc
	done                 chan bool
}

// New creates a new Controller with the given parameters
func New(mac string, pid *pidctrl.PIDController, chamberProvider client.ChamberProvider,
	fermentationProvider client.FermentationProvider, logger log.Logger,
	thermostatOptions ...internal.ThermostatOptionsFunc) *Controller {

	return &Controller{
		ConfigFunc:           ConfigureThermostat,
		mac:                  mac,
		pid:                  pid,
		chamberProvider:      chamberProvider,
		fermentationProvider: fermentationProvider,
		logger:               logger,
		thermostatOptions:    thermostatOptions,
		done:                 make(chan bool, 1),
	}
}

// Start begin the process of polling the service for updates to the chamber
func (c *Controller) Start(interval time.Duration) {

	c.Poll()

	for {
		select {
		case <-c.done:
			return
		case <-time.After(interval):
			c.Poll()
		}
	}
}

// Poll gets a chamber from the server by its mac address. If none is found it create one. It then process the chamber
// to determine what value to set the thermostat to.
func (c *Controller) Poll() {
	var chamber *internal.Chamber
	chamber, err := c.chamberProvider.Get(c.id)
	if err != nil {
		if web.ErrNotFound == errors.Cause(err) {
			c.logger.Println("Chamber does not exist. Creating new chamber")
			chamber := &internal.Chamber{
				Name:       "Chamber " + c.mac,
				MacAddress: c.mac,
				Thermostat: &internal.Thermostat{},
			}
			err := c.chamberProvider.Save(chamber)
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
func (c *Controller) Stop() {
	c.done <- true
}

// processUpdate evaluates the inbound chamber to determine if any changes have occurred
func (c *Controller) processUpdate(chamber *internal.Chamber) error {
	var configChanged bool

	var oldFermID string
	newFermID := chamber.CurrentFermentationID

	if c.Chamber == nil {
		oldFermID = ""
		configChanged = true
		c.Chamber = chamber
	} else {
		oldFermID = c.Chamber.CurrentFermentationID
		configChanged = c.checkChamber(chamber)
	}

	if configChanged {
		c.Chamber.Thermostat.Off()

		c.Chamber = chamber

		err := c.ConfigFunc(c.Chamber.Thermostat, c.pid, c.thermostatOptions...)
		if err != nil {
			return err
		}

		c.Chamber.Thermostat.Subscribe(c.Chamber.ID, c.handleStatusUpdate)
	}

	var err error

	if c.Fermentation == nil {
		if newFermID != "" {
			c.logger.Printf("Fermentation changed from none to %d\n", newFermID)
			c.Fermentation, err = c.getFermentation(newFermID)
			if err != nil {
				return err
			}
		}
	} else {
		if newFermID != "" {
			if oldFermID != newFermID {
				c.logger.Printf("Fermentation changed from %d to %d\n", oldFermID, newFermID)
				c.Chamber.Thermostat.Off()
				c.Fermentation, err = c.getFermentation(newFermID)
				if err != nil {
					return err
				}
			}
		} else {
			c.logger.Printf("Fermentation changed from %d to none\n", oldFermID)
			c.Fermentation = nil
			c.Chamber.Thermostat.Off()
		}
	}

	c.Chamber.CurrentFermentationID = newFermID

	if oldFermID != newFermID || configChanged {
		if c.Fermentation != nil {
			c.logger.Printf("Setting Fermentation to %d\n", c.Fermentation.ID)
			c.Chamber.Thermostat.Set(c.Fermentation.Beer.Schedule[0].TargetTemp)
			if c.Chamber.Thermostat.GetStatus().State == internal.OFF {
				c.Chamber.Thermostat.On()
			}
		}
	}

	return nil
}

func (c *Controller) checkChamber(chamber *internal.Chamber) bool {

	if c.Chamber.Thermostat.ChillerPin != chamber.Thermostat.ChillerPin {
		c.logger.Println("Chiller Pin changed from %s to %s",
			c.Chamber.Thermostat.ChillerPin, chamber.Thermostat.ChillerPin)

		return true
	}

	if c.Chamber.Thermostat.HeaterPin != chamber.Thermostat.HeaterPin {
		c.logger.Println("Heater Pin changed from %s to %s",
			c.Chamber.Thermostat.HeaterPin, chamber.Thermostat.HeaterPin)

		return true
	}

	if c.Chamber.Thermostat.ChillerPin != chamber.Thermostat.ChillerPin {
		c.logger.Println("Thermometer ID changed from %s to %s",
			c.Chamber.Thermostat.ChillerPin, chamber.Thermostat.ChillerPin)

		return true
	}

	return false
}

func (c *Controller) getFermentation(id string) (*internal.Fermentation, error) {

	fermentation, err := c.fermentationProvider.Get(id)
	if err != nil {
		return nil, err
	}

	return fermentation, nil
}

func (c *Controller) handleStatusUpdate(status internal.ThermostatStatus) {

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
		Chamber:     c.Chamber.Name,
		Thermometer: c.Chamber.Thermostat.ThermometerID,
	}

	if c.Fermentation != nil {
		change.FermentationID = c.Fermentation.ID
		change.Beer = c.Fermentation.Beer.Name
	}

	if status.CurrentTemperature != nil {
		change.Temperature = *status.CurrentTemperature
	}

	if err := c.fermentationProvider.SaveTemperatureChange(change); err != nil {
		c.logger.Println(err)
	}
}
