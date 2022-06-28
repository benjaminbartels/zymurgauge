package pid

import (
	"context"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/platform/clock"
	"github.com/benjaminbartels/zymurgauge/internal/temperaturecontrol"
	"github.com/felixge/pidctrl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var _ temperaturecontrol.TemperatureController = (*Controller)(nil)

const (
	pidMin          float64       = 0
	pidMax          float64       = 100
	defaultPeriod   time.Duration = 10 * time.Second
	errorWaitPeriod time.Duration = 10 * time.Second
	dutyTimeDivisor               = 100

	ErrAlreadyRunning   = Error("pid is already running")
	ErrThermometerIsNil = Error("thermometer is nil")
	ErrActuatorIsNil    = Error("actuator is nil")
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Controller struct {
	thermometer device.Thermometer
	actuator    device.Actuator
	pid         *pidctrl.PIDController
	period      time.Duration
	clock       clock.Clock
	logger      *logrus.Logger
	isRunning   bool
	runMutex    sync.Mutex
}

type Status struct {
	Device string
	IsOn   bool
}

func NewController(thermometer device.Thermometer, actuator device.Actuator, kP, kI, kD float64,
	logger *logrus.Logger, options ...OptionsFunc,
) *Controller {
	c := &Controller{
		thermometer: thermometer,

		period:   defaultPeriod,
		actuator: actuator,
		clock:    clock.NewRealClock(),
		logger:   logger,
	}

	pid := pidctrl.NewPIDController(kP, kI, kD)
	pid.SetOutputLimits(pidMin, pidMax)
	c.pid = pid

	for _, option := range options {
		option(c)
	}

	return c
}

type OptionsFunc func(*Controller)

func SetClock(clock clock.Clock) OptionsFunc {
	return func(t *Controller) {
		t.clock = clock
	}
}

func Period(period time.Duration) OptionsFunc {
	return func(t *Controller) {
		t.period = period
	}
}

func (c *Controller) Run(ctx context.Context, setPoint float64) error {
	c.runMutex.Lock()
	if c.isRunning {
		defer c.runMutex.Unlock()

		return ErrAlreadyRunning
	}

	if c.thermometer == nil {
		return ErrThermometerIsNil
	}

	if c.actuator == nil {
		return ErrActuatorIsNil
	}

	c.pid.Set(setPoint)

	c.isRunning = true

	c.runMutex.Unlock()

	lastUpdateTime := c.clock.Now()

	for {
		temperature, err := c.thermometer.GetTemperature()
		if err != nil {
			c.logger.WithError(err).Error("could not read thermometer")
			<-time.After(errorWaitPeriod)

			continue
		}

		c.logger.Debugf("Actuator current temperature is %.4f°C, set point is %.4f°C", temperature, c.pid.Get())

		since := c.clock.Since(lastUpdateTime)
		duty := c.pid.UpdateDuration(temperature, since)

		c.logger.Debugf("Actuator duty is %.2f%%", duty)

		if duty > 0 {
			c.actuator.PWMOn(duty / 100)
		} else {
			c.actuator.Off()
		}

		timer := c.clock.NewTimer(c.period)
		defer timer.Stop()

		select {
		case <-timer.C:
			continue
		case <-ctx.Done():
			c.runMutex.Lock()
			defer c.runMutex.Unlock()
			c.isRunning = false
			if err := c.actuator.Off(); err != nil {
				return errors.Wrap(err, "could not turn actuator off while quiting")
			}

			return nil
		}
	}
}
