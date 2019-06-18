package internal

import (
	"errors"
	"math"
	"reflect"
	"sync"
	"time"
	
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/temporal"
	"github.com/felixge/pidctrl"
)

const (
	outputLimitMin        = -10
	outputLimitMax        = 10
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
	t.status = ThermostatStatus{IsOn: false, State: Idle, CurrentTemperature: nil, Error: nil}
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

func (t *Thermostat) run() {

	lastCheck := t.clock.Now()
	var remainingTime time.Duration

	for t.GetStatus().IsOn {

		if remainingTime == 0 {
			remainingTime = t.interval
		}

		// read temperature
		temperature, err := t.readTemperature()
		if err != nil {
			t.handleError(err)
			return
		}

		t.logf("Temperature is: %.3f.", *temperature)

		t.updateTemperature(temperature)

		var action ThermostatState
		var duration time.Duration

		if remainingTime == t.interval {
			now := t.clock.Now()
			action, duration = t.getNextAction(*temperature, now.Sub(lastCheck))
			lastCheck = now
		} else {
			action = Idle
			duration = remainingTime
		}

		if err := t.performAction(action, duration); err != nil {
			t.handleError(err)
			return
		}

		remainingTime = remainingTime - duration

	}

}

func (t *Thermostat) readTemperature() (*float64, error) {

	// read temperature
	temperature, err := t.thermometer.Read()
	if err != nil {
		return nil, err // ToDo: Wrap error
	}
	if temperature == nil {
		return nil, errors.New("Could not read temperature")
	}
	return temperature, nil

}

func (t *Thermostat) performAction(action ThermostatState, duration time.Duration) error {

	t.logf("%s for %v.", action, duration)

	var err error

	switch action {
	case Idle:
		err = t.idle()
	case Cooling:
		err = t.cool()
	case Heating:
		err = t.heat()
	}

	t.updateState(action) // ToDo: What should state be if error occurred?
	t.sendUpdate()

	if err != nil {
		return err
	}

	select {

	case <-t.quit:
		// Thermostat was set to Off
		if err := t.idle(); err != nil {
			return err
		} else {
			t.updateIsOn(false)
			t.updateState(Idle) // ToDo: What should state be if error occurred?
			t.sendUpdate()
		}

	case <-t.clock.After(duration):

	}

	return nil
}

// getNextAction determines the next action (Heat or Cool or nothing) for the thermostat to perform.  It takes the
// output of the PID controller and converts to to a percentage of the interval.  This percentage is the duration in
// which the thermostat will perform the action
func (t *Thermostat) getNextAction(temperature float64, elapsedTime time.Duration) (ThermostatState, time.Duration) {

	var output float64

	output = t.pid.UpdateDuration(temperature, elapsedTime)

	// get PID max
	_, max := t.pid.OutputLimits()

	// calculate the percentage of interval to be used for the next duration
	percent := math.Abs(math.Round((output/max)/.01) * .01)

	// calculate next duration
	duration := time.Duration(float64(t.interval.Nanoseconds()) * percent)

	action := Idle

	if output < 0 {
		action = Cooling
		if duration > 0 { // ToDo: Implement minChill and minHeat
			if duration < t.minCooling {
				duration = t.minCooling
			}
		}
	} else if output > 0 {
		action = Heating
		if duration > 0 {
			if duration < t.minHeating {
				duration = t.minHeating
			}
		}
	} else {
		duration = t.interval
	}

	t.logf("getNextAction. elapsedTime: %v, in: %.3f, out: %.3f, %.f%%, %s, %v", elapsedTime, temperature, output,
		percent*100, action, duration)

	return action, duration
}

// Off turns the Thermostat Off
func (t *Thermostat) Off() {
	t.quit <- true
}

// Set sets Thermostat to the specified temperature
func (t *Thermostat) Set(temp float64) {
	t.pid.Set(temp)
}

// Subscribe allows a caller to subscribe to ThermostatStatus updates.
func (t *Thermostat) Subscribe(key string, f func(s ThermostatStatus)) {
	t.subs[key] = f
}

// GetStatus return the current status of the Thermostat
func (t *Thermostat) GetStatus() ThermostatStatus {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.status
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

// idle turns the chiller off and the heater off
func (t *Thermostat) idle() error {

	if !reflect.ValueOf(t.chiller).IsNil() {
		if t.status.State == Cooling {
			if err := t.chiller.Off(); err != nil {
				return err
			}
		}
	}

	if !reflect.ValueOf(t.heater).IsNil() {
		if t.status.State == Heating {
			if err := t.heater.Off(); err != nil {
				return err
			}
		}
	}

	return nil
}

// cool turns the chiller on and the heater off
func (t *Thermostat) cool() error {
	if !reflect.ValueOf(t.chiller).IsNil() {
		if t.status.State != Cooling {
			if err := t.chiller.On(); err != nil {
				return err
			}
		}
	}
	if !reflect.ValueOf(t.heater).IsNil() {
		if t.status.State == Heating {
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
		if t.status.State == Cooling {
			if err := t.chiller.Off(); err != nil {
				return err
			}
		}
	}

	if !reflect.ValueOf(t.heater).IsNil() {
		if t.status.State != Heating {
			if err := t.heater.On(); err != nil {
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

// ThermostatState is the current state (Idle, Cooling, Heating) of the thermostat
type ThermostatState string

const (
	// Idle means the Thermostat is not heating or cooling
	Idle ThermostatState = "Idle"
	// Cooling means the Thermostat is cooling
	Cooling ThermostatState = "Cooling"
	// Heating means the Thermostat is heating
	Heating ThermostatState = "Heating"
)
