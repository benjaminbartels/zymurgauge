package internal

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/felixge/pidctrl"
)

type FactoredThermostat struct {
	pidController *pidctrl.PIDController
	thermometer   Thermometer
	chiller       Actuator
	heater        Actuator
	interval      time.Duration
	factor        int
	minChill      time.Duration
	minHeat       time.Duration
	logger        log.Logger
	status        ThermostatStatus
	isOn          AtomBool
	quit          chan bool
	subs          map[string]func(s ThermostatStatus)
	lastCheck     time.Time
	mux           sync.Mutex
}

func NewFactoredThermostat(pidController *pidctrl.PIDController, thermometer Thermometer, chiller, heater Actuator,
	options ...func(*FactoredThermostat) error) (*FactoredThermostat, error) {
	t := &FactoredThermostat{
		pidController: pidController,
		thermometer:   thermometer,
		chiller:       chiller,
		heater:        heater,
		interval:      10 * time.Minute,
		factor:        1,
		minChill:      1 * time.Minute,
		minHeat:       1 * time.Minute,
		quit:          make(chan bool),
		subs:          make(map[string]func(s ThermostatStatus)),
		lastCheck:     time.Now(),
	}

	for _, option := range options {
		err := option(t)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

func Factor(factor int) func(*FactoredThermostat) error {
	return func(t *FactoredThermostat) error {
		t.factor = factor
		return nil
	}
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
	t.pidController.Set(temp)
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

	output := t.pidController.UpdateDuration(temperature, factoredElapsedTime)

	_, max := t.pidController.OutputLimits()

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
