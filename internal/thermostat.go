package internal

import (
	"errors"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/felixge/pidctrl"
)

type Thermostat struct {
	ThermometerID string `json:"thermometerId"`
	ChillerPin    string `json:"chillerPin,omitempty"`
	HeaterPin     string `json:"heaterPin,omitempty"`
	status        ThermostatStatus
	pid           *pidctrl.PIDController
	thermometer   Thermometer
	chiller       Actuator
	heater        Actuator
	minChill      time.Duration
	minHeat       time.Duration
	mux           sync.RWMutex
	quit          chan bool
	subs          map[string]func(s ThermostatStatus)
	logger        log.Logger
}

// ThermostatOptionsFunc is a function that supplies optional configuration to a Thermostat
type ThermostatOptionsFunc func(*Thermostat) error

// MinimumChill set the minimum duration in which the Thermostat will leave the Cooler Actuator On.  This is to prevent
// excessive cycling.  Default is 1 minute.
func MinimumChill(min time.Duration) ThermostatOptionsFunc {
	return func(t *Thermostat) error {
		t.minChill = min
		return nil
	}
}

// MinimumHeat set the minimum duration in which the Thermostat will leave the Heater Actuator On.  This is to prevent
// excessive cycling. Default is 1 minute.
func MinimumHeat(min time.Duration) ThermostatOptionsFunc {
	return func(t *Thermostat) error {
		t.minHeat = min
		return nil
	}
}

// Logger sets the logger to be used.  If not set, nothing is logged.
func Logger(logger log.Logger) ThermostatOptionsFunc {
	return func(t *Thermostat) error {
		t.logger = logger
		return nil
	}
}

// Configure configures a Thermostat with the given parameters
func (t *Thermostat) Configure(pid *pidctrl.PIDController, thermometer Thermometer, chiller, heater Actuator,
	options ...ThermostatOptionsFunc) error {

	pid.SetOutputLimits(-10, 10) // Ensure that limits are set
	t.pid = pid
	t.minChill = 1 * time.Minute
	t.minHeat = 1 * time.Minute
	t.thermometer = thermometer
	t.chiller = chiller
	t.heater = heater
	t.status = ThermostatStatus{IsOn: false, State: IDLE, CurrentTemperature: nil, Error: nil}
	t.mux = sync.RWMutex{}
	t.quit = make(chan bool)
	t.subs = make(map[string]func(s ThermostatStatus))

	for _, option := range options {
		err := option(t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Thermostat) Set(temp float64) {
	t.pid.Set(temp) // ToDo: reset loop?
}

func (t *Thermostat) On() { // ToDo: Refactor this
	t.log("On called")
	t.mux.Lock()
	defer t.mux.Unlock()
	if !t.status.IsOn {
		t.status.IsOn = true
		go t.run()
	}
}

func (t *Thermostat) run() {

	t.mux.RLock()
	isOn := t.status.IsOn
	t.mux.RUnlock()

	for isOn {

		temperature, err := t.thermometer.Read()
		if err != nil {
			t.handleError(err)
			return
		}

		if temperature == nil {
			t.handleError(errors.New("Could not read temperature"))
			return
		}

		output := t.pid.UpdateDuration(*temperature, elapsedTime)

		// ToDo: convert output to duration

		var duration time.Duration

		if output < 0 {

			if duration > 0 {
				if duration < t.minChill {
					duration = t.minChill
				}
			}

			t.coolOn()

		} else if output > 0 {

			if duration > 0 {
				if duration < t.minHeat {
					duration = t.minHeat
				}
			}

			t.heatOn()

		} else {
			duration = 1 * time.Minute
		}

		select {
		case <-t.quit:
			// Thermostat was set to Off
			if err := t.bothOff(); err != nil {
				t.log(err.Error())
				t.handleError(err)
			}
			break
		case <-time.After(duration):
			// Thermostat did work for calculated duration
			t.logf("Acted for %v", duration)
			if err := t.bothOff(); err != nil {
				t.log(err.Error())
				t.handleError(err)
			} else {
				t.updateState(IDLE)
			}
			t.sendUpdate()
		}

	}

	t.sendUpdate()
}

// Off turns the Thermostat Off
func (t *Thermostat) Off() {
	t.log("Off called")
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.status.IsOn {
		t.quit <- true
		t.updateIsOn(false)
	}

}

// GetStatus return the current status of the Thermostat
func (t *Thermostat) GetStatus() ThermostatStatus {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.status
}

// Subscribe allows a caller to subscribe to ThermostatStatus updates.
func (t *Thermostat) Subscribe(key string, f func(s ThermostatStatus)) {
	t.subs[key] = f
}

func (t *Thermostat) coolOn() error {
	return nil
}

func (t *Thermostat) heatOn() error {
	return nil
}

func (t *Thermostat) bothOff() error {
	return nil
}

func (t *Thermostat) updateIsOn(isOn bool) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.status.IsOn = isOn
}

func (t *Thermostat) updateState(state ThermostatState) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.status.State = state
}

func (t *Thermostat) updateTemperature(temperature *float64) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.status.CurrentTemperature = temperature
}

func (t *Thermostat) handleError(err error) {
	t.mux.Lock()
	t.status.IsOn = false
	t.status.Error = err
	t.mux.Unlock()
}

func (t *Thermostat) sendUpdate() {
	t.mux.RLock()
	defer t.mux.RUnlock()
	for _, f := range t.subs {
		if f != nil {
			f(t.status)
		}
	}
}

// log logs the message to the logger if logger has be set
func (t *Thermostat) log(s string) {
	if t.logger != nil {
		t.logger.Println(s)
	}
}

// log logs the formated message to the logger if logger has be set
func (t *Thermostat) logf(s string, a ...interface{}) {
	if t.logger != nil {
		t.logger.Printf(s, a...)
	}
}

// ThermostatStatus contains the state of the thermostat and an error
type ThermostatStatus struct {
	IsOn               bool
	State              ThermostatState
	CurrentTemperature *float64
	Error              error
}

// ThermostatState is the current state (OFF, COOLING, HEATING) of the thermostat
type ThermostatState string

const (
	// IDLE means the Thermostat is not heating or cooling
	IDLE ThermostatState = "IDLE"
	// COOLING means the Thermostat is cooling
	COOLING ThermostatState = "COOLING"
	// HEATING means the Thermostat is heating
	HEATING ThermostatState = "HEATING"
)
