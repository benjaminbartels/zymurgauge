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
	State       State
	target      float64
	isOn        bool
	quit        chan bool
	statusCh    chan Status
	logger      log.Logger
}

// On turns the Thermostat On
func (t *Thermostat) On() chan Status {
	if t.statusCh == nil {
		t.statusCh = make(chan Status)
	}

	if !t.isOn {
		go func() {
			for {
				select {
				case <-t.quit:
					return
				case <-time.After(interval):
					if v, err := t.Thermometer.ReadTemperature(); err != nil {
						t.statusCh <- Status{OFF, err}
					} else if v != nil {
						t.eval(*v)
					}
				}
			}
		}()
	}
	return t.statusCh
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

	var err error

	if temperature > t.target {
		t.log("Temperature above Target. Current: %f Target: %f", temperature, t.target)
		if t.State != COOLING {
			err = t.cool()
		}
	} else if temperature < t.target {
		t.log("Temperature below Target. Current: %f Target: %f", temperature, t.target)
		if t.State != HEATING {
			err = t.heat()
		}
	} else {
		t.log("Temperature equals Target. Current: %f Target: %f", temperature, t.target)
		if t.State != OFF {
			err = t.off()
		}
	}

	if err != nil {
		t.State = ERROR
	}

	t.statusCh <- Status{t.State, err}
}

func (t *Thermostat) cool() error {
	if err := t.Chiller.On(); err != nil {
		return err
	}
	if err := t.Heater.Off(); err != nil {
		return err
	}
	t.State = COOLING
	return nil
}

func (t *Thermostat) heat() error {
	if err := t.Chiller.Off(); err != nil {
		return err
	}
	if err := t.Heater.On(); err != nil {
		return err
	}
	t.State = HEATING
	return nil
}

func (t *Thermostat) off() error {
	if err := t.Chiller.Off(); err != nil {
		return err
	}
	if err := t.Heater.Off(); err != nil {
		return err
	}
	t.State = OFF
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

// Status contains the state of the thermostat and an error
type Status struct {
	State State
	Error error
}

// State is the current state (OFF, COOLING, HEATING) of teh thermostat
type State int

const (
	// OFF means the Thermostat is not heating or cooling
	OFF State = 1 + iota
	// COOLING means the Thermostat is cooling
	COOLING
	// HEATING means the Thermostat is heating
	HEATING
	// ERROR means the Thermostat has an error
	ERROR
)
