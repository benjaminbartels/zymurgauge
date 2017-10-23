package internal

import (
	"fmt"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/ds18b20"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"gobot.io/x/gobot/drivers/gpio"
)

const interval time.Duration = 5 * time.Second

// Thermostat regulates temperature by activating a cooling or heating device when the temperature strays
// from a target
type Thermostat struct {
	ThermometerID string          `json:"thermometerId"`
	ChillerPin    string          `json:"chillerPin,omitempty"`
	HeaterPin     string          `json:"heaterPin,omitempty"`
	State         ThermostatState `json:"-"`
	thermometer   *ds18b20.Thermometer
	chiller       *gpio.RelayDriver
	heater        *gpio.RelayDriver
	target        float64
	isOn          bool
	quit          chan bool
	statusCh      chan ThermostatStatus
	logger        log.Logger
}

// On turns the Thermostat On
func (t *Thermostat) On() chan ThermostatStatus {
	if t.statusCh == nil {
		t.statusCh = make(chan ThermostatStatus)
	}

	if !t.isOn {
		go func() {
			t.checkTemperature()
			for {
				select {
				case <-t.quit:
					return
				case <-time.After(interval):
					t.checkTemperature()
				}
			}
		}()
	}
	return t.statusCh
}

func (t *Thermostat) checkTemperature() {
	if v, err := t.thermometer.ReadTemperature(); err != nil {
		t.statusCh <- ThermostatStatus{ERROR, v, err}
	} else if v != nil {
		t.eval(*v)
	}
}

// Off turns the Thermostat Off
func (t *Thermostat) Off() {
	t.quit <- true
	t.isOn = false
}

// Set sets TemperatureController to the specified temperature
func (t *Thermostat) Set(temp float64) {
	t.log(fmt.Sprintf("Setting thermostat to %f", temp))
	t.target = temp
}

// SetLogger sets the logger
func (t *Thermostat) SetLogger(logger log.Logger) {
	t.logger = logger
}

// InitThermometer initializes the thermometer
func (t *Thermostat) InitThermometer() error {
	thermometer, err := ds18b20.GetThermometer(t.ThermometerID)
	t.thermometer = thermometer
	return err
}

// InitRelays initializes the cooler and heater relays
func (t *Thermostat) InitRelays(w gpio.DigitalWriter) {
	if t.ChillerPin != "" {
		t.chiller = gpio.NewRelayDriver(w, t.ChillerPin)
	}

	if t.HeaterPin != "" {
		t.heater = gpio.NewRelayDriver(w, t.HeaterPin)
	}
}

func (t *Thermostat) eval(temperature float64) {

	var err error

	if temperature > t.target {
		t.log(fmt.Sprintf("Temperature above Target. Current: %f Target: %f", temperature, t.target))
		if t.State != COOLING {
			err = t.cool()
		}
	} else if temperature < t.target {
		t.log(fmt.Sprintf("Temperature below Target. Current: %f Target: %f", temperature, t.target))
		if t.State != HEATING {
			err = t.heat()
		}
	} else {
		t.log(fmt.Sprintf("Temperature equals Target. Current: %f Target: %f", temperature, t.target))
		if t.State != OFF {
			err = t.off()
		}
	}

	if err != nil {
		t.State = ERROR
	}

	t.statusCh <- ThermostatStatus{t.State, &temperature, err}
}

func (t *Thermostat) cool() error {
	if t.chiller != nil {
		if err := t.chiller.On(); err != nil {
			return err
		}
	}
	if t.heater != nil {
		if err := t.heater.Off(); err != nil {
			return err
		}
	}
	t.State = COOLING
	return nil
}

func (t *Thermostat) heat() error {
	if t.chiller != nil {
		if err := t.chiller.Off(); err != nil {
			return err
		}
	}
	if t.heater != nil {
		if err := t.heater.On(); err != nil {
			return err
		}
	}
	t.State = HEATING
	return nil
}

func (t *Thermostat) off() error {
	if t.chiller != nil {
		if err := t.chiller.Off(); err != nil {
			return err
		}
	}
	if t.heater != nil {
		if err := t.heater.Off(); err != nil {
			return err
		}
	}
	t.State = OFF
	return nil
}

func (t *Thermostat) log(s string) {
	if t.logger != nil {
		t.logger.Println(s)
	}
}

// ThermostatStatus contains the state of the thermostat and an error
type ThermostatStatus struct {
	State              ThermostatState
	CurrentTemperature *float64
	Error              error
}

// ThermostatState is the current state (OFF, COOLING, HEATING) of the thermostat
type ThermostatState string

const (
	// OFF means the Thermostat is not heating or cooling
	OFF ThermostatState = "OFF"
	// COOLING means the Thermostat is cooling
	COOLING = "COOLING"
	// HEATING means the Thermostat is heating
	HEATING = "HEATING"
	// ERROR means the Thermostat has an error
	ERROR = "ERROR"
)
