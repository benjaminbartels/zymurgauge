package simulation

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"

	"github.com/benjaminbartels/zymurgauge/internal"
)

// Chamber is used to simulate a real Chamber
type Chamber struct {
	Thermostat              *internal.Thermostat
	beerThermometer         *Thermometer
	chiller                 *Actuator
	heater                  *Actuator
	logger                  log.Logger
	factor                  int
	wallTemp                float64
	airTemp                 float64
	beerTemp                float64
	environmentTemp         float64
	heaterTemp              float64
	beerCapacity            float64
	airCapacity             float64
	wallCapacity            float64
	heaterCapacity          float64
	heaterPower             float64
	coolerPower             float64
	airBeerTransfer         float64
	wallAirTransfer         float64
	heaterAirTransfer       float64
	environmentWallTransfer float64
	heaterToBeer            float64
	heaterToAir             float64
}

// NewChamber creates a new Chamber
func NewChamber(thermostat *internal.Thermostat, beerThermometer *Thermometer,
	chiller, heater *Actuator, factor int, logger log.Logger) *Chamber {
	c := &Chamber{
		Thermostat:      thermostat,
		beerThermometer: beerThermometer,
		chiller:         chiller,
		heater:          heater,
		wallTemp:        20.0,
		airTemp:         20.0,
		beerTemp:        beerThermometer.currentTemp,
		environmentTemp: 20.0,
		heaterTemp:      20.0,
		factor:          factor,
		logger:          logger,
	}

	c.beerCapacity = 4.2 * 1.0 * 20       // heat capacity water * density of water * 20L volume (in kJ per kelvin).
	c.airCapacity = 1.005 * 1.225 * 0.200 // heat capacity of dry air * density of air * 200L volume (in kJ per kelvin).
	// Moist air has only slightly higher heat capacity, 1.02 when saturated at 20C.
	c.wallCapacity = 5.0   // just a guess
	c.heaterCapacity = 1.0 // also a guess, to simulate that heater first heats itself, then starts heating the air

	c.heaterPower = 0.1 // 100W, in kW.
	c.coolerPower = 0.1 // 100W, in kW. Assuming 200W at 50% efficiency

	c.airBeerTransfer = 1.0 / 300
	c.wallAirTransfer = 1.0 / 300
	c.heaterAirTransfer = 1.0 / 30
	c.environmentWallTransfer = 0.001 // losses to environment

	c.heaterToBeer = 0.0 // ratio of heater transferred directly to beer instead of fridge air
	c.heaterToAir = 1.0 - c.heaterToBeer

	c.chiller.Chamber = c
	c.heater.Chamber = c
	c.beerThermometer.Chamber = c

	return c
}

func (c *Chamber) update(onTime time.Duration, t ActuatorType) {

	// fmt.Println("!!!!!! onTime", onTime)

	factoredDuration := onTime * time.Duration(c.factor)

	// fmt.Println("!!!!!! factoredDuration", factoredDuration)

	ticks := int(factoredDuration.Seconds())

	// fmt.Println("!!!!!! updating temp", ticks, "times")

	for i := 0; i < ticks; i++ {

		beerTempNew := c.beerTemp
		airTempNew := c.airTemp
		wallTempNew := c.wallTemp
		heaterTempNew := c.heaterTemp

		beerTempNew += (c.airTemp - c.beerTemp) * c.airBeerTransfer / c.beerCapacity

		if t == Chiller {
			wallTempNew -= c.coolerPower / c.wallCapacity
		} else if t == Heater {
			heaterTempNew += c.heaterPower / c.heaterCapacity
		}

		airTempNew += (c.heaterTemp - c.airTemp) * c.heaterAirTransfer / c.airCapacity
		airTempNew += (c.wallTemp - c.airTemp) * c.wallAirTransfer / c.airCapacity
		airTempNew += (c.beerTemp - c.airTemp) * c.airBeerTransfer / c.airCapacity

		beerTempNew += (c.airTemp - c.beerTemp) * c.airBeerTransfer / c.beerCapacity

		heaterTempNew += (c.airTemp - c.heaterTemp) * c.heaterAirTransfer / c.heaterCapacity

		wallTempNew += (c.environmentTemp - c.wallTemp) * c.environmentWallTransfer / c.wallCapacity
		wallTempNew += (c.airTemp - c.wallTemp) * c.wallAirTransfer / c.wallCapacity

		c.airTemp = airTempNew
		c.beerTemp = beerTempNew
		c.wallTemp = wallTempNew
		c.heaterTemp = heaterTempNew

		//sim yeast fermenting heat
		//c.beerTemp = c.beerTemp + 0.00005

		c.beerThermometer.currentTemp = c.beerTemp

	}
}

// ToDo: remove if not used
// func (c *Chamber) log(s string) {
// 	if c.logger != nil {
// 		c.logger.Println(s)
// 	}
// }

// Thermometer is used to simulate a real Thermometer
type Thermometer struct {
	Chamber     *Chamber
	currentTemp float64
}

// NewThermometer creates a new Thermometer
func NewThermometer(startingTemp float64) *Thermometer {
	t := &Thermometer{
		currentTemp: startingTemp,
	}
	return t
}

// Read returns thr current temperature reading from the Thermometer
func (t *Thermometer) Read() (*float64, error) {

	c := t.Chamber

	if c.chiller.isOn {
		elapsed := time.Since(c.chiller.startTime)
		c.update(elapsed, Chiller)
		c.chiller.startTime = time.Now()
	}

	if c.heater.isOn {
		elapsed := time.Since(c.heater.startTime)
		c.update(elapsed, Heater)
		c.heater.startTime = time.Now()
	}

	return &t.currentTemp, nil
}

// Actuator is used to simulate a real Actuator
type Actuator struct {
	Chamber      *Chamber
	ActuatorType ActuatorType
	isOn         bool
	startTime    time.Time
}

// ActuatorType defines the Actuator's Type
type ActuatorType int

const (
	// Chiller is a type of Actuator
	Chiller ActuatorType = 0
	// Heater is a type of Actuator
	Heater ActuatorType = 1
)

// On turns the Actuator on and records a start time
func (a *Actuator) On() error {

	if !a.isOn {
		a.isOn = true
		a.startTime = time.Now()
	}
	return nil
}

// Off turns the Actuator off, records the elapsed on time and updates the chamber's properties
func (a *Actuator) Off() error {

	if a.isOn {

		elapsed := time.Since(a.startTime)

		a.Chamber.update(elapsed, a.ActuatorType)

		a.isOn = false
	}
	return nil
}