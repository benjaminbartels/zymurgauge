package controller

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
type ChamberCtl struct {
	mac                 string
	chamber             *internal.Chamber
	pid                 *pidctrl.PIDController
	client              *client.Client
	logger              log.Logger
	thermostatOptions   []func(*internal.Thermostat) error
	currentFermentation *internal.Fermentation
	done                chan bool
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

func (c *ChamberCtl) Start() {

	c.poll()

	for {
		select {
		case <-c.done:
			return
		case <-time.After(2 * time.Second):
			c.poll()
		}
	}
}

func (c *ChamberCtl) poll() {
	chamber, err := c.client.ChamberResource.Get(c.mac)
	if err != nil {
		c.logger.Println(err)
	} else {
		if err = c.processUpdate(chamber); err != nil {
			c.logger.Println(err.Error()) // ToDo: Handle
		}
	}
}

func (c *ChamberCtl) Stop() {
	c.done <- true
}

func (c *ChamberCtl) processUpdate(chamber *internal.Chamber) error {
	var err error

	if c.chamber == nil ||
		c.chamber.Thermostat.ChillerPin != chamber.Thermostat.ChillerPin ||
		c.chamber.Thermostat.HeaterPin != chamber.Thermostat.HeaterPin ||
		c.chamber.Thermostat.ThermometerID != chamber.Thermostat.ThermometerID {

		c.logger.Println("New chamber configuration detected")

		if c.chamber != nil {
			c.chamber.Thermostat.Off()
		}

		c.chamber = chamber

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
		err = c.chamber.Thermostat.Configure(c.pid, thermometer, chiller, heater, c.thermostatOptions...)
		if err != nil {
			return err
		}

		c.chamber.Thermostat.Subscribe(c.chamber.MacAddress, c.handleStatusUpdate)
	}

	getFermentation := false

	if c.currentFermentation == nil && c.chamber.CurrentFermentationID != 0 {
		c.logger.Printf("No current fermentation. New fermentation is %d./n",
			c.chamber.CurrentFermentationID, c.chamber.CurrentFermentationID)
		c.chamber.Thermostat.Off()
		getFermentation = true
	} else if c.currentFermentation != nil && c.currentFermentation.ID != c.chamber.CurrentFermentationID {
		c.logger.Printf("Current fermentation is %d. New fermentation is %d./n")
		c.chamber.Thermostat.Off()
		getFermentation = true
	} else if c.currentFermentation != nil && c.currentFermentation.ID == 0 {
		c.logger.Printf("Current fermentation is %d. No fermentation set for chamber./n")
		c.chamber.Thermostat.Off()
		c.currentFermentation = nil
	} else if c.chamber.CurrentFermentationID == 0 {
		c.logger.Println("No fermentation set for chamber")
		c.chamber.Thermostat.Off()
		c.currentFermentation = nil
	}

	if getFermentation {
		c.currentFermentation, err = c.client.FermentationResource.Get(c.chamber.CurrentFermentationID)
		if err != nil {
			return err
		}
	}

	// If there is a fermentation and thermostat is not nil, change setting and turn thermostat On if it is Off
	if c.currentFermentation != nil { // To Do: remove && c.chamber.Thermostat != nil {
		c.chamber.Thermostat.Set(c.currentFermentation.Beer.Schedule[0].TargetTemp)
		if c.chamber.Thermostat.GetStatus().State == internal.OFF {
			c.chamber.Thermostat.On()
		}
	}

	return nil
}

func (c *ChamberCtl) handleStatusUpdate(status internal.ThermostatStatus) {
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
