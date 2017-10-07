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
					if v, err := t.Thermometer.ReadTemperature(); err != nil {
						t.log(err.Error())
						continue
					} else if v != nil {
						if err = t.eval(*v); err != nil {
							t.log(err.Error())
						}
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

func (t *Thermostat) eval(temperature float64) error {

	var err error

	if temperature > t.target {
		t.log("Temperature above Target. Current: %f Target: %f", temperature, t.target)
		if t.state != COOLING {
			err = t.cool()
		}
	} else if temperature < t.target {
		t.log("Temperature below Target. Current: %f Target: %f", temperature, t.target)
		if t.state != HEATING {
			err = t.heat()
		}
	} else {
		t.log("Temperature equals Target. Current: %f Target: %f", temperature, t.target)
		if t.state != OFF {
			err = t.off()
		}
	}

	return err
}

func (t *Thermostat) cool() error {
	if err := t.Chiller.On(); err != nil {
		return err
	}
	if err := t.Heater.Off(); err != nil {
		return err
	}
	t.state = COOLING
	return nil
}

func (t *Thermostat) heat() error {
	if err := t.Chiller.Off(); err != nil {
		return err
	}
	if err := t.Heater.On(); err != nil {
		return err
	}
	t.state = HEATING
	return nil
}

func (t *Thermostat) off() error {
	if err := t.Chiller.Off(); err != nil {
		return err
	}
	if err := t.Heater.Off(); err != nil {
		return err
	}
	t.state = OFF
	return nil
}

func (t *Thermostat) log(s string, a ...interface{}) {
	if t.logger != nil {
		if a != nil {
			t.logger.Printf(s+"\n", a)
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
