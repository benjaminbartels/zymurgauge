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
	pidMin             float64       = 0
	pidMax             float64       = 100
	defaultCyclePeriod time.Duration = 1 * time.Second
	errorWaitPeriod    time.Duration = 10 * time.Second
	dutyTimeDivisor                  = 100

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
	cyclePeriod time.Duration
	clock       clock.Clock
	logger      *logrus.Logger
	isRunning   bool
	runMutex    sync.Mutex
}

type Status struct {
	Device string
	IsOn   bool
}

func NewPIDTemperatureController(thermometer device.Thermometer, actuator device.Actuator, kP, kI, kD float64,
	logger *logrus.Logger, options ...OptionsFunc,
) *Controller {
	c := &Controller{
		thermometer: thermometer,

		cyclePeriod: defaultCyclePeriod,
		actuator:    actuator,
		clock:       clock.NewRealClock(),
		logger:      logger,
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

// CyclePeriod sets the duration of the chiller's PWM cycle.
func CyclePeriod(period time.Duration) OptionsFunc {
	return func(t *Controller) {
		t.cyclePeriod = period
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

	return c.startCycle(ctx)
}

func (c *Controller) startCycle(ctx context.Context) error {
	lastUpdateTime := c.clock.Now()

	for {
		temperature, err := c.thermometer.GetTemperature()
		if err != nil {
			c.logger.WithError(err).Error("could not read thermometer")
			<-time.After(errorWaitPeriod)

			continue
		}

		since := c.clock.Since(lastUpdateTime)
		dutyCycle := c.pid.UpdateDuration(temperature, since)
		lastUpdateTime = c.clock.Now()
		dutyTime := time.Duration(float64(c.cyclePeriod.Nanoseconds()) * dutyCycle / dutyTimeDivisor)
		waitTime := c.cyclePeriod - dutyTime

		c.logger.Debugf("Actuator current temperature is %.4f°C, set point is %.4f°C", temperature, c.pid.Get())
		c.logger.Debugf("Actuator dutyCycle is %.2f%%, dutyTime is %s, waitTime is %s",
			dutyCycle, dutyTime, waitTime)

		if dutyTime > 0 {
			if err := c.actuator.On(); err != nil {
				c.logger.WithError(err).Error("could not turn actuator on")
				<-time.After(errorWaitPeriod)

				return nil
			}

			c.logger.Debugf("Actuator acting for %v", dutyTime)

			if didComplete := c.wait(ctx, dutyTime); !didComplete {
				return c.quit(c.actuator)
			}

			c.logger.Debugf("Actuator acted for %v", dutyTime)
		}

		if waitTime > 0 {
			if err := c.actuator.Off(); err != nil {
				c.logger.WithError(err).Error("could not turn actuator off")
				<-time.After(errorWaitPeriod)

				return nil
			}

			c.logger.Debugf("Actuator waiting for %v", waitTime)

			if didComplete := c.wait(ctx, waitTime); !didComplete {
				return c.quit(c.actuator)
			}

			c.logger.Debugf("Actuator waited for %v", waitTime)
		}
	}
}

func (c *Controller) wait(ctx context.Context, waitTime time.Duration) bool {
	timer := c.clock.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		c.runMutex.Lock()
		defer c.runMutex.Unlock()
		c.isRunning = false

		return false
	}
}

func (c *Controller) quit(actuator device.Actuator) error {
	c.logger.Debug("Actuator quiting")

	if err := actuator.Off(); err != nil {
		return errors.Wrap(err, "could not turn actuator off while quiting")
	}

	return nil
}
