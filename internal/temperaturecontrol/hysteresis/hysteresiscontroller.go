package hysteresis

import (
	"context"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/temperaturecontrol"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var _ temperaturecontrol.TemperatureController = (*Controller)(nil)

const (
	defaultCyclePeriod     = 10 * time.Second
	defaultChillerCooldown = 10 * time.Minute
	errorWaitPeriod        = 10 * time.Second

	ErrAlreadyRunning   = Error("pid is already running")
	ErrThermometerIsNil = Error("thermometer is nil")
	ErrActuatorIsNil    = Error("actuator is nil")
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Controller struct {
	thermometer         device.Thermometer
	chiller             device.Actuator
	heater              device.Actuator
	hysteresisBand      float64
	cyclePeriod         time.Duration
	chillerCooldown     time.Duration
	logger              *logrus.Logger
	setPoint            float64
	isRunning           bool
	chillerOnStartTime  time.Time
	chillerOffStartTime time.Time
	setPointCh          chan float64
	runMutex            sync.Mutex
}

func NewController(thermometer device.Thermometer, chiller, heater device.Actuator,
	hysteresisBand float64,
	logger *logrus.Logger, options ...OptionsFunc) *Controller {
	t := &Controller{
		thermometer:     thermometer,
		chiller:         chiller,
		heater:          heater,
		hysteresisBand:  hysteresisBand,
		cyclePeriod:     defaultCyclePeriod,
		chillerCooldown: defaultChillerCooldown,
		setPointCh:      make(chan float64),
		logger:          logger,
	}

	for _, option := range options {
		option(t)
	}

	return t
}

type OptionsFunc func(*Controller)

func CyclePeriod(cyclePeriod time.Duration) OptionsFunc {
	return func(t *Controller) {
		t.cyclePeriod = cyclePeriod
	}
}

func ChillerCooldown(chillerCooldown time.Duration) OptionsFunc {
	return func(t *Controller) {
		t.chillerCooldown = chillerCooldown
	}
}

func (c *Controller) SetTemperature(temperature float64) {
	c.setPointCh <- temperature
}

func (c *Controller) Run(ctx context.Context, setPoint float64) error {
	c.runMutex.Lock()
	if c.isRunning {
		defer c.runMutex.Unlock()

		return ErrAlreadyRunning
	}

	if c.thermometer == nil {
		defer c.runMutex.Unlock()

		return ErrThermometerIsNil
	}

	if c.chiller == nil || c.heater == nil {
		defer c.runMutex.Unlock()

		return ErrActuatorIsNil
	}

	c.setPoint = setPoint

	c.isRunning = true

	c.runMutex.Unlock()

	if err := c.startCycle(ctx); err != nil {
		return errors.Wrap(err, "could not start cycle")
	}

	return nil
}

func (c *Controller) startCycle(ctx context.Context) error {
	for {
		var err error

		temperature, err := c.thermometer.GetTemperature()
		if err != nil {
			c.logger.WithError(err).Error("could not read thermometer")
			<-time.After(errorWaitPeriod)

			continue
		}

		divisor := 2.0
		upperBound := c.setPoint + c.hysteresisBand/divisor
		lowerBound := c.setPoint - c.hysteresisBand/divisor

		if temperature > upperBound {
			err = c.startChilling()
		} else if temperature < lowerBound {
			err = c.startHeating()
		}

		if err != nil {
			c.logger.WithError(err).Error(err, "errors occurred while evaluating temperature")
			<-time.After(errorWaitPeriod)

			continue
		}

		if didComplete := c.wait(ctx); !didComplete {
			return c.quit()
		}
	}
}

func (c *Controller) startChilling() error {
	var result error

	cooldownOverTime := c.chillerOffStartTime.Add(c.chillerCooldown)

	if time.Now().After(cooldownOverTime) {
		if err := c.chiller.On(); err != nil {
			result = multierror.Append(result, errors.Wrap(err, "could not turn chiller actuator on"))
		} else {
			c.chillerOnStartTime = time.Now()
		}
	} else {
		c.logger.Debugf("Cannot turn chiller on for another %s", time.Until(cooldownOverTime))
	}

	if err := c.heater.Off(); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "could not turn heater actuator off"))
	}

	if result != nil {
		return errors.Wrap(result, "errors occurred while starting to chill")
	}

	return nil
}

func (c *Controller) startHeating() error {
	var result error

	err := c.chiller.Off()
	if err != nil {
		result = multierror.Append(result, errors.Wrap(err, "could not turn chiller actuator off"))
	} else if !c.chillerOnStartTime.IsZero() {
		c.chillerOffStartTime = time.Now()
	}

	if heaterErr := c.heater.On(); err != nil {
		result = multierror.Append(result, errors.Wrap(heaterErr, "could not turn heater actuator on"))
	}

	if result != nil {
		return errors.Wrap(result, "errors occurred while starting to heat")
	}

	return nil
}

func (c *Controller) wait(ctx context.Context) bool {
	timer := time.NewTimer(c.cyclePeriod)
	defer timer.Stop()

	select {
	case <-timer.C:
		return true
	case temperature := <-c.setPointCh:
		c.setPoint = temperature

		return true
	case <-ctx.Done():
		c.runMutex.Lock()
		defer c.runMutex.Unlock()
		c.isRunning = false

		return false
	}
}

func (c *Controller) quit() error {
	var result error

	if err := c.chiller.Off(); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "could not turn chiller actuator off"))
	}

	if err := c.heater.Off(); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "could not turn heater actuator off"))
	}

	if result != nil {
		return errors.Wrap(result, "error(s) occurred while quitting")
	}

	return nil
}
