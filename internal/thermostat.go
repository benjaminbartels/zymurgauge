package internal

import (
	"time"

	"gobot.io/x/gobot/drivers/gpio"

	"github.com/benjaminbartels/zymurgauge/internal/ds18b20"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
)

const interval time.Duration = time.Minute

// Thermostat regulates temperature by activating a cooling or heating device when the temperature strays
// from a target
type Thermostat struct {
	Thermometer *ds18b20.Thermometer `json:"thermometer"`
	Chiller     *gpio.RelayDriver    `json:"chiller"`
	Heater      *gpio.RelayDriver    `json:"heater"`
	target      float64
	isOn        bool
	state       state
	quit        chan bool
	logger      log.Logger
}

// On turns the Thermostat On
func (t *Thermostat) On() {
	if !t.isOn {
		go func() {
			for {
				select {
				case <-t.quit:
					return
				case <-time.After(interval):
					v, err := t.Thermometer.ReadTemperature()
					if err != nil {
						if t.logger != nil {
							t.logger.Println(err)
						}
						continue
					}
					if v != nil {
						t.eval(*v)
					}
				}
			}
		}()
	}
}

// Off turns the Thermostat Off
func (t *Thermostat) Off() {
	t.quit <- true
	t.isOn = false
}

// Set sets TemperatureController to the specified temperature
func (t *Thermostat) Set(temp float64) {
	t.target = temp
}

func (t *Thermostat) eval(temperature float64) {
	if temperature > t.target {
		t.log("Temperature above Target. Current: %f Target: %f\n", temperature, t.target)
		if t.state != COOLING {
			t.Chiller.On()
			t.Heater.Off()
			t.state = COOLING
		}
	} else if temperature < t.target {
		t.log("Temperature below Target. Current: %f Target: %f\n", temperature, t.target)
		if t.state != HEATING {
			t.Chiller.Off()
			t.Heater.On()
			t.state = HEATING
		}
	} else {
		t.log("Temperature equals Target. Current: %f Target: %f\n", temperature, t.target)
		if t.state != OFF {
			t.Chiller.Off()
			t.Heater.Off()
			t.state = OFF
		}
	}
}

func (t *Thermostat) log(s string, a ...interface{}) {
	if t.logger != nil {
		if a != nil {
			t.logger.Printf(s, a)
		} else {
			t.logger.Println(s)
		}
	}
}

type state int

const (
	// OFF means the Thermostat is not heating or cooling
	OFF state = 1 + iota
	// COOLING means the Thermostat is cooling
	COOLING
	// HEATING means the Thermostat is heating
	HEATING
)
