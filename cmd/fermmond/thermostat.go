package main

import (
	"errors"
	"math"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/gpio"
	"github.com/benjaminbartels/zymurgauge/internal/platform/atomic"
	"github.com/benjaminbartels/zymurgauge/internal/platform/ds18b20"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/felixge/pidctrl"
)

// ToDo: Determine in Factored thermostat is needed.

// Thermostat regulates the temperature of Chamber by monitoring the temperature that is read from a Thermometer and
// switching on and off a heater and cooler Actuators.  It determines when to switching Actuators on and off by feeding
// inputs into a PIDController and reading the outputs
type Thermostat struct {
	pid         *pidctrl.PIDController
	thermometer *ds18b20.Thermometer
	chiller     *gpio.Actuator
	heater      *gpio.Actuator
	interval    time.Duration
	minChill    time.Duration
	minHeat     time.Duration
	logger      log.Logger
	status      ThermostatStatus
	isOn        atomic.Bool
	quit        chan bool
	subs        map[string]func(s ThermostatStatus)
	lastCheck   time.Time
}

// NewThermostat creates a new Thermostat with the given parameters
func NewThermostat(thermostatDefinition internal.Thermostat, pid *pidctrl.PIDController, logger log.Logger,
	options ...func(*Thermostat) error) (*Thermostat, error) {

	var err error

	t := &Thermostat{
		pid:       pid,
		logger:    logger,
		interval:  10 * time.Minute,
		minChill:  1 * time.Minute,
		minHeat:   1 * time.Minute,
		quit:      make(chan bool),
		subs:      make(map[string]func(s ThermostatStatus)),
		lastCheck: time.Now(),
	}

	t.thermometer, err = ds18b20.NewThermometer(thermostatDefinition.ThermometerID)
	if err != nil {
		return nil, err
	}

	t.chiller, err = gpio.NewActuator(thermostatDefinition.ChillerPin)
	if err != nil {
		return nil, err
	}

	t.heater, err = gpio.NewActuator(thermostatDefinition.HeaterPin)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		err := option(t)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

// On turns the Thermostat on and allows to being monitoring
func (t *Thermostat) On() { // ToDo: Refactor this
	if !t.isOn.Get() {
		go func() {
			t.isOn.Set(true)

			for {

				var action ThermostatState
				var duration time.Duration

				temperature, err := t.thermometer.Read()
				if temperature == nil {
					err = errors.New("Could not read temperature")
				} else {
					action, duration = t.getNextAction(*temperature)
					if action == COOLING {
						err = t.cool()
					} else if action == HEATING {
						err = t.heat()
					}

				}

				if err != nil {
					t.log(err.Error())
					t.updateStatus(ERROR, nil, err)
					if err := t.off(); err != nil {
						t.log(err.Error())
					}
					return
				}

				t.updateStatus(action, temperature, nil)

				select {
				case <-t.quit:
					if err := t.off(); err != nil {
						t.updateStatus(ERROR, nil, err)
						t.log(err.Error())
					} else {
						t.updateStatus(OFF, temperature, nil)
					}
					return
				case <-time.After(t.interval):
					t.logf("Acted for entire interval of %v", t.interval)
				case <-time.After(duration):
					t.logf("Acted for only %v", duration)
					if err := t.off(); err != nil {
						t.updateStatus(ERROR, nil, err)
						t.log(err.Error())
					} else {
						t.updateStatus(OFF, temperature, nil)
					}
					t.logf("Waiting for remaining interval time %v", t.interval-duration)
					<-time.After(t.interval - duration)
				}
			}
		}()
	}
}

// Off turns the Thermostat Off
func (t *Thermostat) Off() {
	t.quit <- true
	t.isOn.Set(false)
}

// Set sets TemperatureController to the specified temperature
func (t *Thermostat) Set(temp float64) {
	t.pid.Set(temp)
}

// GetStatus return the current status of the Thermostat
func (t *Thermostat) GetStatus() ThermostatStatus {
	return t.status
}

// Subscribe allows a caller to subscribe to ThermostatStatus updates.
func (t *Thermostat) Subscribe(key string, f func(s ThermostatStatus)) {
	t.subs[key] = f
}

//
func (t *Thermostat) getNextAction(temperature float64) (ThermostatState, time.Duration) {

	now := time.Now()

	elapsedTime := now.Sub(t.lastCheck)

	output := t.pid.UpdateDuration(temperature, elapsedTime)

	_, max := t.pid.OutputLimits()

	percent := math.Abs(math.Round((output/max)/.01) * .01)

	duration := time.Duration(float64(t.interval.Nanoseconds()) * percent)

	action := OFF

	if output < 0 {
		action = COOLING
		// if duration > 0 {
		// 	if duration < t.minChill {
		// 		duration = t.minChill
		// 	}
		// }
	} else if output > 0 {
		action = HEATING
		// if duration > 0 {
		// 	if duration < t.minHeat {
		// 		duration = t.minHeat
		// 	}
		// }
	} else {
		duration = t.interval
	}

	t.logf("getNextAction - lastCheck: %v, elapsedTime: %v, in: %.3f, out: %.3f, %.f%%, %s, %v",
		t.lastCheck.Format("2006/01/02 15:04:05"), elapsedTime, temperature, output, percent*100,
		action, duration)

	t.lastCheck = now

	return action, duration
}

func (t *Thermostat) updateStatus(state ThermostatState, temperature *float64, err error) {
	t.status = ThermostatStatus{
		State:              state,
		CurrentTemperature: temperature,
		Error:              err,
	}

	for _, f := range t.subs {
		if f != nil {
			f(t.status)
		}
	}
}

func (t *Thermostat) cool() error {
	if t.chiller != nil {
		if t.status.State != COOLING {
			if err := t.chiller.On(); err != nil {
				return err
			}
		}
	}
	if t.heater != nil {
		if t.status.State == HEATING {
			if err := t.heater.Off(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Thermostat) heat() error {
	if t.chiller != nil {
		if t.status.State == COOLING {
			if err := t.chiller.Off(); err != nil {
				return err
			}
		}
	}
	if t.heater != nil {
		if t.status.State != HEATING {
			if err := t.heater.On(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Thermostat) off() error {

	if t.chiller != nil {
		if t.status.State == COOLING {
			if err := t.chiller.Off(); err != nil {
				return err
			}
		}
	}
	if t.heater != nil {
		if t.status.State == HEATING {
			if err := t.heater.Off(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Thermostat) log(s string) {
	if t.logger != nil {
		t.logger.Println(s)
	}
}

func (t *Thermostat) logf(s string, a ...interface{}) {
	//t.log(fmt.Sprintf(s, a))
	if t.logger != nil {
		t.logger.Printf(s, a...)
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
	COOLING ThermostatState = "COOLING"
	// HEATING means the Thermostat is heating
	HEATING ThermostatState = "HEATING"
	// ERROR means the Thermostat has an error
	ERROR ThermostatState = "ERROR"
)

// Interval sets the interval in which the Thermostat checks the temperature.  Default is 10 minutes.
func Interval(interval time.Duration) func(*Thermostat) error {
	return func(t *Thermostat) error {
		t.interval = interval
		return nil
	}
}

// MinimumChill set the minimum duration in which the Thermostat will leave the Cooler Actuator On.  This is to prevent
// excessive cycling.  Default is 1 minute.
func MinimumChill(min time.Duration) func(*Thermostat) error {
	return func(t *Thermostat) error {
		t.minChill = min
		return nil
	}
}

// MinimumHeat set the minimum duration in which the Thermostat will leave the Heater Actuator On.  This is to prevent
// excessive cycling. Default is 1 minute.
func MinimumHeat(min time.Duration) func(*Thermostat) error {
	return func(t *Thermostat) error {
		t.minHeat = min
		return nil
	}
}
