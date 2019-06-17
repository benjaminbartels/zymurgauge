package internal

import (
	"errors"
	"math"
	"reflect"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/temporal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/felixge/pidctrl"
)

const (
	outputLimitMin = -1
	outputLimitMax = 10
	defaultMinimumCooling = 1 * time.Minute
	defaultMinimumHeating = 1 * time.Minute
)

// Thermostat regulates the temperature of Chamber by monitoring the temperature that is read from a Thermometer and
// switching on and off a heater and cooler Actuators.  It determines when to switching Actuators on and off by feeding
// inputs into a PIDController and reading the outputs
type Thermostat struct {
	ThermometerID string `json:"thermometerId"`
	ChillerPin    string `json:"chillerPin,omitempty"`
	HeaterPin     string `json:"heaterPin,omitempty"`
	pid           *pidctrl.PIDController
	thermometer   Thermometer
	chiller       Actuator
	heater        Actuator
	clock         temporal.Clock
	interval      time.Duration
	minCooling    time.Duration
	minHeating    time.Duration
	logger        log.Logger
	status        ThermostatStatus
	mux           sync.RWMutex
	quit          chan bool
	subs          map[string]func(s ThermostatStatus)
}

// ThermostatOptionsFunc is a function that supplies optional configuration to a Thermostat
type ThermostatOptionsFunc func(*Thermostat) error

// Logger sets the logger to be used.  If not set, nothing is logged.
func Clock(clock temporal.Clock) ThermostatOptionsFunc {
	return func(t *Thermostat) error {
		t.clock = clock
		return nil
	}
}

// Interval sets the interval in which the Thermostat checks the temperature.  Default is 10 minutes.
func Interval(interval time.Duration) ThermostatOptionsFunc {
	return func(t *Thermostat) error {
		t.interval = interval
		return nil
	}
}

// MinimumCooling set the minimum duration in which the Thermostat will leave the Chiller Actuator On.  This is to prevent
// excessive cycling.  Default is 1 minute.
func MinimumCooling(min time.Duration) ThermostatOptionsFunc {
	return func(t *Thermostat) error {
		t.minCooling = min
		return nil
	}
}

// MinimumHeating set the minimum duration in which the Thermostat will leave the Heater Actuator On.  This is to prevent
// excessive cycling. Default is 1 minute.
func MinimumHeating(min time.Duration) ThermostatOptionsFunc {
	return func(t *Thermostat) error {
		t.minHeating = min
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

	pid.SetOutputLimits(outputLimitMin, outputLimitMax) // Ensure that limits are set
	t.pid = pid
	t.interval = 10 * time.Minute
	t.minCooling = defaultMinimumCooling
	t.minHeating = defaultMinimumHeating
	t.quit = make(chan bool)
	t.subs = make(map[string]func(s ThermostatStatus))
	t.thermometer = thermometer
	t.clock = temporal.NewClock()
	t.chiller = chiller
	t.heater = heater
	t.status = ThermostatStatus{IsOn: false, State: IDLING, CurrentTemperature: nil, Error: nil}
	t.mux = sync.RWMutex{}

	for _, option := range options {
		err := option(t)
		if err != nil {
			return err
		}
	}

	return nil
}

// On turns the Thermostat on and allows to being monitoring
func (t *Thermostat) On() { // ToDo: Refactor this
	t.log("On called")
	t.mux.Lock()
	defer t.mux.Unlock()
	if !t.status.IsOn {
		t.status.IsOn = true
		go t.run()
	}
}

func (t *Thermostat) readTemperature() (float64, error) {

		// read temperature
		temperature, err := t.thermometer.Read()
		if err != nil {
			return 0, err // ToDo: Wrap error
		}
		if temperature == nil {
			return 0, errors.New("Could not read temperature")
		}
		return *temperature, nil


}

func (t *Thermostat) run() {

	lastCheck := time.Now()

	for t.GetStatus().IsOn {

		now := t.clock.Now()

		t.logf("Stating new cycle @ %s", now)

		// read temperature
		temperature, err := t.readTemperature()
		if err != nil {
			t.handleError(err)
			return
		}

		elapsedTime := now.Sub(lastCheck)

		action, duration := t.getNextAction(temperature, elapsedTime)
		
		if err := t.performAction(); err != nil {
			t.handleError(err)
			return
		}

	}

}

func (t *Thermostat) performAction(state ThermostatState, duration time.Duration) error {

	var err error

	switch state {
	case IDLING:
		err = t.idle()
	case COOLING:
		err = t.cool()
	case HEATING:
		err = t.heat()
	}

	t.updateState(state) // ToDo: What should state be if error occurred?

	if err != nil {
		return err
	}

	select {

	case <-t.quit:
		// Thermostat was set to Off
		if err := t.idle(); err != nil {
			return err
		} else {
			t.updateState(IDLING) // ToDo: What should state be if error occurred?
		}

	case <-t.clock.After(duration):
		// Thermostat did work for calculated duration, turn off until next interval
		t.logf("Acted for only %v", duration) // ToDo: better log message
	}


	return nil
}

// wait waits for a the thermostat to be turned Off OR the interval to elapse OR or the calculated duration to
// elapse
func (t *Thermostat) wait(temperature *float64, duration time.Duration) {

	select {
	case <-t.quit:
		// Thermostat was set to Off
		if err := t.idle(); err != nil {
			t.handleError(err)
		} else {
			t.updateState(IDLING)
		}
		return
	case <-t.clock.After(t.interval):
		// Thermostat did work for the entire duration
		t.logf("Acted for entire interval of %v", t.interval)
	case <-t.clock.After(duration):
		// Thermostat did work for calculated duration, turn off until next interval
		t.logf("Acted for only %v", duration)
		if err := t.idle(); err != nil {
			t.handleError(err)
		} else {
			t.updateState(IDLING)
		}
		t.logf("Waiting for remaining interval time %v", t.interval-duration)
		<-t.clock.After(t.interval - duration)
	}
}

// getNextAction determines the next action (Heat or Cool or nothing) for the thermostat to perform.  It takes the
// output of the PID controller and converts to to a percentage of the interval.  This percentage is the duration in
// which the thermostat will perform the action
func (t *Thermostat) getNextAction(temperature float64, elapsedTime time.Duration) (ThermostatState, time.Duration) {

	output := t.pid.UpdateDuration(temperature, elapsedTime)
	
	// get PID max
	_, max := t.pid.OutputLimits()

	// calculate the percentage of interval to be used for the next duration
	percent := math.Abs(math.Round((output/max)/.01) * .01)

	// calculate next duration
	duration := time.Duration(float64(t.interval.Nanoseconds()) * percent)

	action := IDLING

	if output < 0 {
		action = COOLING
		if duration > 0 { // ToDo: Implement minChill and minHeat
			if duration < t.minCooling {
				duration = t.minCooling
			}
		}
	} else if output > 0 {
		action = HEATING
		if duration > 0 {
			if duration < t.minHeating {
				duration = t.minHeating
			}
		}
	} else {
		duration = t.interval
	}

	t.logf("getNextAction - lastCheck: %v, elapsedTime: %v, in: %.3f, out: %.3f, %.f%%, %s, %v",
		t.lastCheck.Format("2006/01/02 15:04:05"), elapsedTime, temperature, output, percent*100,
		action, duration)

	return action, duration
}

// Off turns the Thermostat Off
func (t *Thermostat) Off() {
	t.log("Off called")
	t.updateIsOn(false)
	t.quit <- true
}

// Set sets Thermostat to the specified temperature
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
	t.log(err.Error())
	t.sendUpdate()

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

// cool turns the chiller on and the heater off
func (t *Thermostat) cool() error {
	if !reflect.ValueOf(t.chiller).IsNil() {
		if t.status.State != COOLING {
			if err := t.chiller.On(); err != nil {
				return err
			}
		}
	}
	if !reflect.ValueOf(t.heater).IsNil() {
		if t.status.State == HEATING {
			if err := t.heater.Off(); err != nil {
				return err
			}
		}
	}

	return nil
}

// heat turns the chiller off and the heater on
func (t *Thermostat) heat() error {
	if !reflect.ValueOf(t.chiller).IsNil() {
		if t.status.State == COOLING {
			if err := t.chiller.Off(); err != nil {
				return err
			}
		}
	}

	if !reflect.ValueOf(t.heater).IsNil() {
		if t.status.State != HEATING {
			if err := t.heater.On(); err != nil {
				return err
			}
		}
	}

	return nil
}

// idle turns the chiller off and the heater off
func (t *Thermostat) idle() error {

	if !reflect.ValueOf(t.chiller).IsNil() {
		if t.status.State == COOLING {
			if err := t.chiller.Off(); err != nil {
				return err
			}
		}
	}

	if !reflect.ValueOf(t.heater).IsNil() {
		if t.status.State == HEATING {
			if err := t.heater.Off(); err != nil {
				return err
			}
		}
	}

	return nil
}

// log logs the message to the logger if logger has be set
func (t *Thermostat) log(s string) {
	if t.logger != nil {
		t.logger.Println(s)
	}
}

// log logs the formated message to the logger if logger has be set
func (t *Thermostat) logf(s string, a ...interface{}) {
	//t.log(fmt.Sprintf(s, a))
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

// ThermostatState is the current state (IDLING, COOLING, HEATING) of the thermostat
type ThermostatState string

const (
	// IDLING means the Thermostat is not heating or cooling
	IDLING ThermostatState = "IDLING"
	// COOLING means the Thermostat is cooling
	COOLING ThermostatState = "COOLING"
	// HEATING means the Thermostat is heating
	HEATING ThermostatState = "HEATING"
)
