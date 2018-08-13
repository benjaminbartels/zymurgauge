package simulation

import (
	"errors"
	"math"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/atomic"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/felixge/pidctrl"
)

type FactoredThermostat struct {
	ThermometerID string `json:"thermometerId"`
	ChillerPin    string `json:"chillerPin,omitempty"`
	HeaterPin     string `json:"heaterPin,omitempty"`
	pid           *pidctrl.PIDController
	thermometer   internal.Thermometer
	chiller       internal.Actuator
	heater        internal.Actuator
	interval      time.Duration
	factor        int
	minChill      time.Duration
	minHeat       time.Duration
	logger        log.Logger
	status        ThermostatStatus
	isOn          atomic.Bool
	quit          chan bool
	subs          map[string]func(s ThermostatStatus)
	lastCheck     time.Time
}

// func NewFactoredThermostat(pid *pidctrl.PIDController, thermometer Thermometer, chiller, heater Actuator,
// 	options ...func(*FactoredThermostat) error) (*FactoredThermostat, error) {
// 	t := &FactoredThermostat{
// 		pid:         pid,
// 		thermometer: thermometer,
// 		chiller:     chiller,
// 		heater:      heater,
// 		interval:    10 * time.Minute,
// 		factor:      1,
// 		minChill:    1 * time.Minute,
// 		minHeat:     1 * time.Minute,
// 		quit:        make(chan bool),
// 		subs:        make(map[string]func(s ThermostatStatus)),
// 		lastCheck:   time.Now(),
// 	}

// 	for _, option := range options {
// 		err := option(t)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return t, nil
// }

// Interval sets the interval in which the Thermostat checks the temperature.  Default is 10 minutes.
func Interval(interval time.Duration) func(*FactoredThermostat) error {
	return func(t *FactoredThermostat) error {
		t.interval = interval
		return nil
	}
}

// MinimumChill set the minimum duration in which the Thermostat will leave the Cooler Actuator On.  This is to prevent
// excessive cycling.  Default is 1 minute.
func MinimumChill(min time.Duration) func(*FactoredThermostat) error {
	return func(t *FactoredThermostat) error {
		t.minChill = min
		return nil
	}
}

// MinimumHeat set the minimum duration in which the Thermostat will leave the Heater Actuator On.  This is to prevent
// excessive cycling. Default is 1 minute.
func MinimumHeat(min time.Duration) func(*FactoredThermostat) error {
	return func(t *FactoredThermostat) error {
		t.minHeat = min
		return nil
	}
}

// Logger sets the logger to be used.  If not set, nothing is logged.
func Logger(logger log.Logger) func(*FactoredThermostat) error {
	return func(t *FactoredThermostat) error {
		t.logger = logger
		return nil
	}
}

func Factor(factor int) func(*FactoredThermostat) error {
	return func(t *FactoredThermostat) error {
		t.factor = factor
		return nil
	}
}

// Configure configures a Thermostat with the given parameters
func (t *FactoredThermostat) Configure(pid *pidctrl.PIDController, thermometer internal.Thermometer,
	chiller, heater internal.Actuator, options ...func(*FactoredThermostat) error) error {

	t.pid = pid
	t.interval = 10 * time.Minute
	t.factor = 1
	t.minChill = 1 * time.Minute
	t.minHeat = 1 * time.Minute
	t.quit = make(chan bool)
	t.subs = make(map[string]func(s ThermostatStatus))
	t.lastCheck = time.Now()
	t.thermometer = thermometer
	t.chiller = chiller
	t.heater = heater

	for _, option := range options {
		err := option(t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *FactoredThermostat) On() {
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
				} else {
					t.updateStatus(action, temperature, nil)
				}

				select {
				case <-t.quit:
					return
				case <-time.After(t.interval):
					t.logf("Acted for entire interval of %v", t.applyFactor(t.interval))
				case <-time.After(duration):
					t.logf("Acted for only %v", t.applyFactor(duration))
					if err := t.off(); err != nil {
						t.updateStatus(ERROR, nil, err)
						t.log(err.Error())
					} else {
						t.updateStatus(OFF, temperature, nil)
					}
					t.logf("Waiting for remaining interval time %v", t.applyFactor(t.interval-duration))
					<-time.After(t.interval - duration)
				}
			}
		}()
	}
}

// Off turns the Thermostat Off
func (t *FactoredThermostat) Off() {
	t.quit <- true
	t.isOn.Set(false)
}

// Set sets TemperatureController to the specified temperature
func (t *FactoredThermostat) Set(temp float64) {
	t.pid.Set(temp)
}

func (t *FactoredThermostat) GetStatus() ThermostatStatus {
	return t.status
}

func (t *FactoredThermostat) Subscribe(key string, f func(s ThermostatStatus)) {
	t.subs[key] = f
}

func (t *FactoredThermostat) getNextAction(temperature float64) (ThermostatState, time.Duration) {

	now := time.Now()

	elapsedTime := now.Sub(t.lastCheck)

	factoredElapsedTime := elapsedTime * (time.Duration(t.factor))

	output := t.pid.UpdateDuration(temperature, factoredElapsedTime)

	_, max := t.pid.OutputLimits()

	percent := math.Abs(math.Round((output/max)/.01) * .01)

	//duration := time.Duration(math.Abs(percent*float64(t.interval.Seconds()))) * time.Second
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

	t.logf("getNextAction - lastCheck: %v, elapsedTime: %v, factoredElapsedTime: %v, in: %.3f, out: %.3f, %.f%%, %s, %v",
		t.lastCheck.Format("2006/01/02 15:04:05"), elapsedTime, factoredElapsedTime, temperature, output, percent*100,
		action, t.applyFactor(duration))

	t.lastCheck = now

	return action, duration
}

func (t *FactoredThermostat) updateStatus(state ThermostatState, temperature *float64, err error) {
	t.status = ThermostatStatus{state, temperature, err}

	for _, f := range t.subs {
		if f != nil {
			f(t.status)
		}
	}
}

func (t *FactoredThermostat) cool() error {
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

func (t *FactoredThermostat) heat() error {
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

func (t *FactoredThermostat) off() error {

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

func (t *FactoredThermostat) applyFactor(d time.Duration) time.Duration {
	return d * (time.Duration(t.factor))
}

func (t *FactoredThermostat) log(s string) {
	if t.logger != nil {
		t.logger.Println(s)
	}
}

func (t *FactoredThermostat) logf(s string, a ...interface{}) {
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
