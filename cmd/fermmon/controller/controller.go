//nolint:funlen,gocognit,nestif
package controller

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/client"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
	"github.com/benjaminbartels/zymurgauge/internal/thermostat"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Controller controls a Chamber. ConfigFunc is the function used to configure the controller's Thermostat.
type Controller struct {
	CreateFunc func(string, string, string, float64, float64, float64, float64, float64, float64, *logrus.Logger,
		...thermostat.OptionsFunc) (*thermostat.Thermostat, error)
	Chamber              *storage.Chamber
	Thermostat           *thermostat.Thermostat
	Fermentation         *storage.Fermentation
	mac                  string
	chillerKp            float64
	chillerKi            float64
	chillerKd            float64
	heaterKp             float64
	heaterKi             float64
	heaterKd             float64
	chamberProvider      client.ChamberProvider
	fermentationProvider client.FermentationProvider
	logger               *logrus.Logger
	thermostatOptions    []thermostat.OptionsFunc
	done                 chan bool
}

// New creates a new Controller with the given parameters.
func New(mac string, chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	chamberProvider client.ChamberProvider, fermentationProvider client.FermentationProvider, logger *logrus.Logger,
	thermostatOptions ...thermostat.OptionsFunc) *Controller {
	return &Controller{
		CreateFunc:           CreateThermostat,
		mac:                  mac,
		chillerKp:            chillerKp,
		chillerKi:            chillerKi,
		chillerKd:            chillerKd,
		heaterKp:             heaterKp,
		heaterKi:             heaterKi,
		heaterKd:             heaterKd,
		chamberProvider:      chamberProvider,
		fermentationProvider: fermentationProvider,
		logger:               logger,
		done:                 make(chan bool, 1),
	}
}

// Start begin the process of polling the service for updates to the chamber.
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
	var chamber *storage.Chamber

	chamber, err := c.chamberProvider.Get(c.mac)

	if errors.Is(err, web.ErrNotFound) {
		c.logger.Println("Chamber does not exist. Creating new chamber")

		chamber := &storage.Chamber{
			Name:       "Chamber " + c.mac,
			MacAddress: c.mac,
			Thermostat: &storage.Thermostat{},
		}
		if err := c.chamberProvider.Save(chamber); err != nil {
			c.logger.Fatalln(err.Error()) // ToDo: Handle
		}
	} else if err = c.processUpdate(chamber); err != nil {
		c.logger.Println(err.Error()) // ToDo: Handle
	}
}

// Stop ends the polling process.
func (c *Controller) Stop() {
	c.done <- true
}

// processUpdate evaluates the inbound chamber to determine if any changes have occurred.
func (c *Controller) processUpdate(chamber *storage.Chamber) error {
	var (
		configChanged bool
		oldFermID     uint64
	)

	newFermID := chamber.CurrentFermentationID

	if c.Chamber == nil {
		oldFermID = 0
		configChanged = true
		c.Chamber = chamber
	} else {
		oldFermID = c.Chamber.CurrentFermentationID
		configChanged = c.checkChamber(chamber)
	}

	if configChanged {
		c.Thermostat.Off()

		c.Chamber = chamber

		newThermostat, err := c.CreateFunc(c.Chamber.Thermostat.ThermometerID,
			c.Chamber.Thermostat.ChillerPin, c.Chamber.Thermostat.HeaterPin,
			c.chillerKp, c.chillerKi, c.chillerKd, c.heaterKp, c.heaterKi, c.heaterKd, c.logger, c.thermostatOptions...)
		if err != nil {
			return err
		}

		c.Thermostat = newThermostat
	}

	var err error

	if c.Fermentation != nil {
		if newFermID != 0 {
			c.logger.Printf("Fermentation changed from none to %d\n", newFermID)

			c.Fermentation, err = c.getFermentation(newFermID)
			if err != nil {
				return err
			}
		}
	} else {
		if newFermID != 0 {
			if oldFermID != newFermID {
				c.logger.Printf("Fermentation changed from %d to %d\n", oldFermID, newFermID)
				c.Thermostat.Off()
				c.Fermentation, err = c.getFermentation(newFermID)
				if err != nil {
					return err
				}
			}
		} else {
			c.logger.Printf("Fermentation changed from %d to none\n", oldFermID)
			c.Fermentation = nil
			c.Thermostat.Off()
		}
	}

	c.Chamber.CurrentFermentationID = newFermID

	if oldFermID != newFermID || configChanged {
		if c.Fermentation != nil {
			c.logger.Printf("Setting Fermentation to %d\n", c.Fermentation.ID)

			go func() {
				if err := c.Thermostat.On(c.Fermentation.Beer.Schedule[0].TargetTemp); err != nil {
					c.logger.WithError(err).Warn("Thermostat failed to turn on")
				}
			}()
		}
	}

	return nil
}

func (c *Controller) checkChamber(chamber *storage.Chamber) bool {
	if c.Chamber.Thermostat.ChillerPin != chamber.Thermostat.ChillerPin {
		c.logger.Printf("Chiller Pin changed from %s to %s\n",
			c.Chamber.Thermostat.ChillerPin, chamber.Thermostat.ChillerPin)

		return true
	}

	if c.Chamber.Thermostat.HeaterPin != chamber.Thermostat.HeaterPin {
		c.logger.Printf("Heater Pin changed from %s to %s\n",
			c.Chamber.Thermostat.HeaterPin, chamber.Thermostat.HeaterPin)

		return true
	}

	if c.Chamber.Thermostat.ChillerPin != chamber.Thermostat.ChillerPin {
		c.logger.Printf("Thermometer ID changed from %s to %s\n",
			c.Chamber.Thermostat.ChillerPin, chamber.Thermostat.ChillerPin)

		return true
	}

	return false
}

func (c *Controller) getFermentation(id uint64) (*storage.Fermentation, error) {
	fermentation, err := c.fermentationProvider.Get(id)
	if err != nil {
		return nil, err
	}

	return fermentation, nil
}
