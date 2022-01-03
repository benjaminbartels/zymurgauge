package pid

import (
	"context"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/platform/clock"
	"github.com/felixge/pidctrl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	pidMin                     float64       = 0
	pidMax                     float64       = 100
	defaultChillingCyclePeriod time.Duration = 30 * time.Minute
	defaultHeatingCyclePeriod  time.Duration = 10 * time.Minute
	defaultChillingMinimum     time.Duration = 10 * time.Minute
	defaultHeatingMinimum      time.Duration = 10 * time.Second
	dutyCycleMultiplyer                      = 100
)

var ErrAlreadyRunning = errors.New("pid is already running")

type TemperatureController struct {
	thermometer         device.Thermometer
	chiller             device.Actuator
	heater              device.Actuator
	chillerKp           float64
	chillerKi           float64
	chillerKd           float64
	heaterKp            float64
	heaterKi            float64
	heaterKd            float64
	chillingCyclePeriod time.Duration
	heatingCyclePeriod  time.Duration
	chillingMinimum     time.Duration
	heatingMinimum      time.Duration
	clock               clock.Clock
	logger              *logrus.Logger
	isRunning           bool
	runMutex            sync.Mutex
}

func NewPIDTemperatureController(thermometer device.Thermometer, chiller, heater device.Actuator,
	chillerKp, chillerKi, chillerKd, heaterKp, heaterKi, heaterKd float64,
	logger *logrus.Logger, options ...OptionsFunc) *TemperatureController {
	t := &TemperatureController{
		thermometer:         thermometer,
		chillerKp:           chillerKp,
		chillerKi:           chillerKi,
		chillerKd:           chillerKd,
		heaterKp:            heaterKp,
		heaterKi:            heaterKi,
		heaterKd:            heaterKd,
		chiller:             chiller,
		chillingCyclePeriod: defaultChillingCyclePeriod,
		heatingCyclePeriod:  defaultHeatingCyclePeriod,
		chillingMinimum:     defaultChillingMinimum,
		heatingMinimum:      defaultHeatingMinimum,
		clock:               clock.NewRealClock(),
		heater:              heater,
		logger:              logger,
	}

	for _, option := range options {
		option(t)
	}

	return t
}

type OptionsFunc func(*TemperatureController)

func SetClock(clock clock.Clock) OptionsFunc {
	return func(t *TemperatureController) {
		t.clock = clock
	}
}

// ChillingCyclePeriod sets the duration of the chiller's PWM cycle.  Default is 30 minutes.
func SetChillingCyclePeriod(period time.Duration) OptionsFunc {
	return func(t *TemperatureController) {
		t.chillingCyclePeriod = period
	}
}

// HeatingCyclePeriod sets the duration of the chiller's PWM cycle.  Default is 1 minute.
func SetHeatingCyclePeriod(period time.Duration) OptionsFunc {
	return func(t *TemperatureController) {
		t.heatingCyclePeriod = period
	}
}

// SetChillingMinimum sets the minimum duration in which the Temperature Controller will leave the Chiller Actuator On.
// This is to prevent excessive cycling.  Default is 10 minutes.
func SetChillingMinimum(min time.Duration) OptionsFunc {
	return func(t *TemperatureController) {
		t.chillingMinimum = min
	}
}

// SetHeaterMinimum sets the minimum duration in which the Temperature Controller will leave the Heater Actuator On.
// This is to prevent excessive cycling. Default is 10 seconds.
func SetHeatingMinimum(min time.Duration) OptionsFunc {
	return func(t *TemperatureController) {
		t.heatingMinimum = min
	}
}

func (t *TemperatureController) startCycle(ctx context.Context, name string, pid *pidctrl.PIDController,
	actuator device.Actuator, period, minimum time.Duration) error {
	lastUpdateTime := t.clock.Now()

	for {
		temperature, err := t.thermometer.GetTemperature()
		if err != nil {
			return errors.Wrap(err, "could not read thermometer")
		}

		since := t.clock.Since(lastUpdateTime)
		output := pid.UpdateDuration(temperature, since)
		dutyCycle := output / pidMax
		lastUpdateTime = t.clock.Now()
		dutyTime := time.Duration(float64(period.Nanoseconds()) * dutyCycle)
		waitTime := period - dutyTime

		t.logger.Debugf("Actuator %s set point is %.4f°C", name, pid.Get())
		t.logger.Debugf("Actuator %s current temperature is %.4f°C", name, temperature)

		p, i, d := pid.PID()
		t.logger.Debugf("Actuator %s PID is %f, %f, %f", name, p, i, d)
		t.logger.Debugf("Actuator %s time since last update is %s", name, since)
		t.logger.Debugf("Actuator %s output is %f", name, output)
		t.logger.Debugf("Actuator %s dutyCycle is %.2f%%", name, dutyCycle*dutyCycleMultiplyer)
		t.logger.Debugf("Actuator %s dutyTime is %s", name, dutyTime)
		t.logger.Debugf("Actuator %s waitTime is %s", name, waitTime)

		if dutyTime > 0 {
			if dutyTime < minimum {
				t.logger.Debugf("Forcing %s actuator to a run for a minimum of %s", name, minimum)
				dutyTime = minimum
			}

			if err := actuator.On(); err != nil {
				return errors.Wrapf(err, "could not turn %s actuator on", name)
			}

			t.logger.Debugf("Actuator %s acting for %v", name, dutyTime)

			if didComplete := t.wait(ctx, dutyTime); !didComplete {
				return t.quit(name, actuator)
			}

			t.logger.Debugf("Actuator %s acted for %v", name, dutyTime)
		}

		if waitTime > 0 {
			if err := actuator.Off(); err != nil {
				return errors.Wrap(err, "could not turn actuator off")
			}

			t.logger.Debugf("Actuator %s waiting for %v", name, waitTime)

			if didComplete := t.wait(ctx, waitTime); !didComplete {
				return t.quit(name, actuator)
			}

			t.logger.Debugf("Actuator %s waited for %v", name, waitTime)
		}
	}
}

func (t *TemperatureController) wait(ctx context.Context, waitTime time.Duration) bool {
	timer := t.clock.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		t.runMutex.Lock()
		defer t.runMutex.Unlock()
		t.isRunning = false

		return false
	}
}

func (t *TemperatureController) Run(ctx context.Context, setPoint float64) error {
	t.runMutex.Lock()
	if t.isRunning {
		defer t.runMutex.Unlock()

		return ErrAlreadyRunning
	}

	t.isRunning = true

	t.runMutex.Unlock()

	chillerPID := newPID(t.chillerKp, t.chillerKi, t.chillerKd, pidMin, pidMax)
	chillerPID.Set(setPoint)

	heaterPID := newPID(t.heaterKp, t.heaterKi, t.heaterKd, pidMin, pidMax)
	heaterPID.Set(setPoint)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return t.startCycle(ctx, "chiller", chillerPID, t.chiller, t.chillingCyclePeriod, t.chillingMinimum)
	})
	g.Go(func() error {
		return t.startCycle(ctx, "heater", heaterPID, t.heater, t.heatingCyclePeriod, t.heatingMinimum)
	})

	if err := g.Wait(); err != nil {
		return errors.Wrap(err, "failure while waiting")
	}

	return nil
}

func newPID(kP, kI, kD, min, max float64) *pidctrl.PIDController {
	pid := pidctrl.NewPIDController(kP, kI, kD)
	pid.SetOutputLimits(min, max)

	return pid
}

func (t *TemperatureController) quit(name string, actuator device.Actuator) error {
	t.logger.Debugf("Actuator %s quiting", name)

	if err := actuator.Off(); err != nil {
		return errors.Wrap(err, "could not turn actuator off while quiting")
	}

	return nil
}
