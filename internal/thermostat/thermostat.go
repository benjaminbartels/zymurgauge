package thermostat

import (
	"context"
	"sync"
	"time"

	"github.com/felixge/pidctrl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// Thermometer represents a device that and read temperatures.
type Thermometer interface {
	Read() (float64, error)
}

// Actuator represents a device that can be switched on and off.
type Actuator interface {
	On() error
	Off() error
}

const (
	pidMin                    float64       = 0
	pidMax                    float64       = 100
	defaultChillerCyclePeriod time.Duration = 30 * time.Minute
	defaultHeatingCyclePeriod time.Duration = 11 * time.Minute
	defaultChillingMinimum    time.Duration = 10 * time.Minute
	defaultHeatingMinimum     time.Duration = 10 * time.Second
)

var ErrAlreadyOn = errors.New("thermostat is already on")

type OptionsFunc func(*Thermostat)

func SetClock(clock Clock) OptionsFunc {
	return func(t *Thermostat) {
		t.clock = clock
	}
}

// ChillerCyclePeriod sets the duration of the chiller's PWM cycle.  Default is 30 minutes.
func SetChillerCyclePeriod(period time.Duration) OptionsFunc {
	return func(t *Thermostat) {
		t.chillerCyclePeriod = period
	}
}

// HeatingCyclePeriod sets the duration of the chiller's PWM cycle.  Default is 1 minute.
func SetHeatingCyclePeriod(period time.Duration) OptionsFunc {
	return func(t *Thermostat) {
		t.heatingCyclePeriod = period
	}
}

// SetChillingMinimum sets the minimum duration in which the Thermostat will leave the Chiller Actuator On.
// This is to prevent excessive cycling.  Default is 10 minutes.
func SetChillingMinimum(min time.Duration) OptionsFunc {
	return func(t *Thermostat) {
		t.chillingMinimum = min
	}
}

// SetHeaterMinimum sets the minimum duration in which the Thermostat will leave the Heater Actuator On.
// This is to prevent excessive cycling. Default is 10 seconds.
func SetHeatingMinimum(min time.Duration) OptionsFunc {
	return func(t *Thermostat) {
		t.heatingMinimum = min
	}
}

type Thermostat struct {
	thermometer        Thermometer
	chiller            Actuator
	heater             Actuator
	chillerKp          float64
	chillerKi          float64
	chillerKd          float64
	heaterKp           float64
	heaterKi           float64
	heaterKd           float64
	chillerCyclePeriod time.Duration
	chillingMinimum    time.Duration
	heatingCyclePeriod time.Duration
	heatingMinimum     time.Duration
	clock              Clock
	logger             *logrus.Logger
	isOn               bool
	onMutex            sync.Mutex
	cancelFn           context.CancelFunc
}

func NewThermostat(thermometer Thermometer, chiller, heater Actuator,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...OptionsFunc) *Thermostat {
	t := &Thermostat{
		thermometer:        thermometer,
		chillerKp:          chillerKp,
		chillerKi:          chillerKi,
		chillerKd:          chillerKd,
		heaterKp:           heaterKp,
		heaterKi:           heaterKi,
		heaterKd:           heaterKd,
		chiller:            chiller,
		chillerCyclePeriod: defaultChillerCyclePeriod,
		heatingCyclePeriod: defaultHeatingCyclePeriod,
		chillingMinimum:    defaultChillingMinimum,
		heatingMinimum:     defaultHeatingMinimum,
		clock:              NewRealClock(),
		heater:             heater,
		logger:             logger,
	}

	for _, option := range options {
		option(t)
	}

	return t
}

func (t *Thermostat) startCycle(ctx context.Context, name string, pid *pidctrl.PIDController,
	actuator Actuator, period, minimum time.Duration) error {
	lastUpdateTime := t.clock.Now()

	for {
		temperature, err := t.thermometer.Read()
		if err != nil {
			return errors.Wrap(err, "could not read thermometer")
		}

		output := pid.UpdateDuration(temperature, time.Since(lastUpdateTime))
		dutyCycle := output / pidMax
		lastUpdateTime = t.clock.Now()

		dutyTime := time.Duration(float64(period.Nanoseconds()) * dutyCycle)

		waitTime := period

		if dutyTime > 0 {
			if dutyTime < minimum {
				t.logger.Infof("Forcing %s actuator to a run for a minimum of %s", name, minimum)
				dutyTime = minimum
			}

			if err := actuator.On(); err != nil {
				return errors.Wrapf(err, "could not turn %s actuator on", name)
			}

			dutyTimer := time.NewTimer(dutyTime)

			t.logger.Infof("Actuator %s acting for %v", name, dutyTime)

			select {
			case <-dutyTimer.C:
				t.logger.Infof("Actuator %s acted for %v", name, dutyTime)

				if err := actuator.Off(); err != nil {
					return errors.Wrap(err, "could not turn actuator off after duty cycle")
				}

			case <-ctx.Done():
				return t.quit(dutyTimer, actuator)
			}

			waitTime -= dutyTime
		} else {
			t.logger.Infof("Actuator %s waiting for %v", name, waitTime)
		}

		waitTimer := time.NewTimer(waitTime)

		select {
		case <-waitTimer.C:
			t.logger.Infof("Actuator %s waited for %v", name, waitTime)
		case <-ctx.Done():
			return t.quit(waitTimer, actuator)
		}
	}
}

func (t *Thermostat) On(setPoint float64) error {
	t.onMutex.Lock()
	if t.isOn {
		defer t.onMutex.Unlock()
		return ErrAlreadyOn
	}

	t.isOn = true
	t.onMutex.Unlock()

	chillerPID := newPID(t.chillerKp, t.chillerKi, t.chillerKd, pidMin, pidMax)
	chillerPID.Set(setPoint)

	heaterPID := newPID(t.heaterKp, t.heaterKi, t.heaterKd, pidMin, pidMax)
	heaterPID.Set(setPoint)

	ctx, cancel := context.WithCancel(context.Background())
	t.cancelFn = cancel

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return t.startCycle(ctx, "chiller", chillerPID, t.chiller, t.chillerCyclePeriod, t.chillingMinimum)
	})
	g.Go(func() error {
		return t.startCycle(ctx, "heater", heaterPID, t.heater, t.heatingCyclePeriod, t.heatingMinimum)
	})

	return g.Wait()
}

func (t *Thermostat) Off() {
	t.cancelFn()
}

func newPID(kP, kI, kD, min, max float64) *pidctrl.PIDController {
	pid := pidctrl.NewPIDController(kP, kI, kD)
	pid.SetOutputLimits(min, max)

	return pid
}

func (t *Thermostat) quit(timer *time.Timer, actuator Actuator) error {
	timer.Stop()

	if err := actuator.Off(); err != nil {
		return errors.Wrap(err, "could not turn actuator off while quiting")
	}

	return nil
}
