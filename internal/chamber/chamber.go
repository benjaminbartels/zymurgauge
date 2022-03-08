package chamber

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/pid"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
	"github.com/benjaminbartels/zymurgauge/internal/platform/metrics"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Chamber represents an insulated box (fridge) with internal heating/cooling elements that reacts to changes in
// monitored temperatures, by correcting small deviations from your desired fermentation temperature.
type Chamber struct {
	ID                      string                  `json:"id"` // TODO: omit empty?
	Name                    string                  `json:"name"`
	DeviceConfig            DeviceConfig            `json:"deviceConfig"` // TODO: make is a struct and not a list
	ChillerKp               float64                 `json:"chillerKp"`
	ChillerKi               float64                 `json:"chillerKi"`
	ChillerKd               float64                 `json:"chillerKd"`
	HeaterKp                float64                 `json:"heaterKp"`
	HeaterKi                float64                 `json:"heaterKi"`
	HeaterKd                float64                 `json:"heaterKd"`
	CurrentBatch            *brewfather.BatchDetail `json:"currentBatch,omitempty"`
	ModTime                 time.Time               `json:"modTime"`
	CurrentFermentationStep *string                 `json:"currentFermentationStep,omitempty"`
	Readings                *Readings               `json:"readings,omitempty"`
	logger                  *logrus.Logger
	metrics                 metrics.Metrics
	beerThermometer         device.Thermometer
	auxiliaryThermometer    device.Thermometer
	externalThermometer     device.Thermometer
	hydrometer              device.Hydrometer
	chiller                 device.Actuator
	heater                  device.Actuator
	temperatureController   device.TemperatureController
	service                 brewfather.Service
	cancelFunc              context.CancelFunc
	metricsInterval         time.Duration
	runMutex                *sync.Mutex
	readingsMutex           *sync.RWMutex
}

type DeviceConfig struct {
	ChillerGPIO              string `json:"chillerGpio"`
	HeaterGPIO               string `json:"heaterGpio"`
	BeerThermometerType      string `json:"beerThermometerType"`
	BeerThermometerID        string `json:"beerThermometerId"`
	AuxiliaryThermometerType string `json:"auxiliaryThermometerType"`
	AuxiliaryThermometerID   string `json:"auxiliaryThermometerId"`
	ExternalThermometerType  string `json:"externalThermometerType"`
	ExternalThermometerID    string `json:"externalThermometerId"`
	HydrometerType           string `json:"hydrometerType"`
	HydrometerID             string `json:"hydrometerId"`
}

type Readings struct {
	BeerTemperature      *float64 `json:"beerTemperature,omitempty"`
	AuxiliaryTemperature *float64 `json:"auxiliaryTemperature,omitempty"`
	ExternalTemperature  *float64 `json:"externalTemperature,omitempty"`
	HydrometerGravity    *float64 `json:"hydrometerGravity,omitempty"`
}

// TODO: refactor to use generics in the future.

func (c *Chamber) Configure(configurator Configurator, service brewfather.Service,
	logger *logrus.Logger, metrics metrics.Metrics, metricsInterval time.Duration) error {
	c.service = service
	c.logger = logger
	c.metrics = metrics
	c.metricsInterval = metricsInterval

	errs := c.configureDevices(configurator, c.DeviceConfig)

	c.temperatureController = pid.NewPIDTemperatureController(
		c.beerThermometer, c.chiller, c.heater, c.ChillerKp, c.ChillerKi, c.ChillerKd,
		c.HeaterKp, c.HeaterKi, c.HeaterKd, logger)

	c.runMutex = &sync.Mutex{}
	c.readingsMutex = &sync.RWMutex{}

	if errs != nil {
		return &InvalidConfigurationError{configErrors: errs}
	}

	return nil
}

func (c *Chamber) configureDevices(configurator Configurator, config DeviceConfig) []error {
	var errs []error

	errs = append(errs, c.configureActuators(configurator, config)...)

	errs = append(errs, c.configureThermometers(configurator, config)...)

	if config.HydrometerType != "" {
		h, err := getHydrometer(configurator, config.HydrometerType,
			config.HydrometerID)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "could not configure hydrometer"))
		}

		c.hydrometer = h
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func (c *Chamber) configureActuators(configurator Configurator, config DeviceConfig) []error {
	var errs []error

	a, err := configurator.CreateGPIOActuator(config.ChillerGPIO)
	if err != nil {
		errs = append(errs, errors.Wrapf(err, "could not create new GPIO %s for chiller", config.ChillerGPIO))
	}

	c.chiller = a

	a, err = configurator.CreateGPIOActuator(config.HeaterGPIO)
	if err != nil {
		errs = append(errs, errors.Wrapf(err, "could not create new GPIO %s for heater", config.HeaterGPIO))
	}

	c.heater = a

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func (c *Chamber) configureThermometers(configurator Configurator, config DeviceConfig) []error {
	var errs []error

	if config.BeerThermometerType != "" {
		t, err := getThermometer(configurator, config.BeerThermometerType,
			config.BeerThermometerID)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "could not configure beer thermometer"))
		}

		c.beerThermometer = t
	}

	if config.AuxiliaryThermometerType != "" {
		t, err := getThermometer(configurator, config.AuxiliaryThermometerType,
			config.AuxiliaryThermometerID)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "could not configure auxiliary thermometer"))
		}

		c.auxiliaryThermometer = t
	}

	if config.ExternalThermometerType != "" {
		t, err := getThermometer(configurator, config.ExternalThermometerType,
			config.ExternalThermometerID)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "could not configure external thermometer"))
		}

		c.externalThermometer = t
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func getThermometer(configurator Configurator, thermometerType, id string) (device.Thermometer, error) {
	switch thermometerType {
	case "ds18b20": // TODO: make enum
		createdDevice, err := configurator.CreateDs18b20(id)
		if err != nil {
			return nil, errors.Wrapf(err, "could not create new Ds18b20 %s", id)
		}

		return createdDevice, nil
	case "tilt":
		createdDevice, err := configurator.CreateTilt(tilt.Color(id))
		if err != nil {
			return nil, errors.Wrapf(err, "could not create new %s Tilt", id)
		}

		return createdDevice, nil
	default:
		return nil, errors.Errorf("invalid thermometer type '%s'", thermometerType)
	}
}

func getHydrometer(configurator Configurator, hydrometerType, id string) (device.Hydrometer, error) {
	switch hydrometerType {
	case "tilt":
		createdDevice, err := configurator.CreateTilt(tilt.Color(id))
		if err != nil {
			return nil, errors.Wrapf(err, "could not create new %s Tilt", id)
		}

		return createdDevice, nil
	default:
		return nil, errors.Errorf("invalid hydrometer type '%s'", hydrometerType)
	}
}

func (c *Chamber) StartFermentation(ctx context.Context, stepID string) error {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if c.CurrentBatch == nil {
		return ErrNoCurrentBatch
	}

	step := c.getStep(stepID)
	if step == nil {
		return ErrInvalidStep
	}

	if c.cancelFunc != nil {
		c.cancelFunc()
	}

	temp := step.StepTemperature
	ctx, cancelFunc := context.WithCancel(ctx)
	c.cancelFunc = cancelFunc

	go func() {
		c.updateReadings(ctx)

		for {
			timer := time.NewTimer(c.metricsInterval)
			defer timer.Stop()

			select {
			case <-timer.C:
				c.updateReadings(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	// TODO: this should return an error to the called, but would require the manager to keep track of go routines
	// The manager would keep a map of cancelFunc, that would be returned from chamber.StartFermentation() along
	// with an error.
	go func() {
		if err := c.temperatureController.Run(ctx, temp); err != nil {
			c.logger.WithError(err).Errorf("could not run temperature controller for chamber %s", c.Name)
			c.cancelFunc = nil // TODO: test this
		}
	}()

	return nil
}

func (c *Chamber) StopFermentation() error {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if c.cancelFunc == nil {
		return ErrNotFermenting
	}

	c.cancelFunc()

	c.cancelFunc = nil

	return nil
}

func (c *Chamber) IsFermenting() bool {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	return c.cancelFunc != nil
}

func (c *Chamber) GetReadings() Readings {
	c.readingsMutex.RLock()
	defer c.readingsMutex.RUnlock()

	return *c.Readings
}

func (c *Chamber) updateReadings(ctx context.Context) {
	c.readingsMutex.Lock()
	defer c.readingsMutex.Unlock()

	c.Readings = &Readings{}

	var v *float64

	var err error

	if v, err = c.getBeerTemperature(); err != nil {
		if !errors.Is(err, ErrDeviceIsNil) {
			c.logger.WithError(err).Error("could not get reading for beer temperature")
		}
	}

	c.Readings.BeerTemperature = v

	if v, err = c.getAuxiliaryTemperature(); err != nil {
		if !errors.Is(err, ErrDeviceIsNil) {
			c.logger.WithError(err).Error("could not get reading for auxiliary temperature")
		}
	}

	c.Readings.AuxiliaryTemperature = v

	if v, err = c.getExternalTemperature(); err != nil {
		if !errors.Is(err, ErrDeviceIsNil) {
			c.logger.WithError(err).Error("could not get reading for external temperature")
		}
	}

	c.Readings.ExternalTemperature = v

	if v, err = c.getHydrometerGravity(); err != nil {
		if !errors.Is(err, ErrDeviceIsNil) {
			c.logger.WithError(err).Error("could not get reading for hydrometer gravity")
		}
	}

	c.Readings.HydrometerGravity = v

	c.emitMetrics()

	if err := c.sendToBrewFather(ctx); err != nil {
		c.logger.WithError(err).Error("Unable to send readings to Brewfather")
	}
}

func (c *Chamber) getBeerTemperature() (*float64, error) {
	if c.beerThermometer == nil {
		return nil, ErrDeviceIsNil
	}

	t, err := c.beerThermometer.GetTemperature()
	if err != nil {
		return nil, errors.Wrap(err, "could not get beer temperature")
	}

	return &t, nil
}

func (c *Chamber) getAuxiliaryTemperature() (*float64, error) {
	if c.auxiliaryThermometer == nil {
		return nil, ErrDeviceIsNil
	}

	t, err := c.auxiliaryThermometer.GetTemperature()
	if err != nil {
		return nil, errors.Wrap(err, "could not get auxiliary temperature")
	}

	return &t, nil
}

func (c *Chamber) getExternalTemperature() (*float64, error) {
	if c.externalThermometer == nil {
		return nil, ErrDeviceIsNil
	}

	t, err := c.externalThermometer.GetTemperature()
	if err != nil {
		return nil, errors.Wrap(err, "could not get external temperature")
	}

	return &t, nil
}

func (c *Chamber) getHydrometerGravity() (*float64, error) {
	if c.hydrometer == nil {
		return nil, ErrDeviceIsNil
	}

	t, err := c.hydrometer.GetGravity()
	if err != nil {
		return nil, errors.Wrap(err, "could not get hydrometer gravity")
	}

	return &t, nil
}

func (c *Chamber) getStep(stepID string) *brewfather.FermentationStep {
	var step *brewfather.FermentationStep

	for i := range c.CurrentBatch.Recipe.Fermentation.Steps {
		if c.CurrentBatch.Recipe.Fermentation.Steps[i].Type == stepID {
			step = &c.CurrentBatch.Recipe.Fermentation.Steps[i]

			break
		}
	}

	return step
}

func (c *Chamber) emitMetrics() {
	if c.Readings.BeerTemperature != nil {
		c.metrics.Gauge(fmt.Sprintf("zymurgauge.%s.beer_temperature,sensor_id=%s", c.Name,
			c.beerThermometer.GetID()), *c.Readings.BeerTemperature)
	}

	if c.Readings.AuxiliaryTemperature != nil {
		c.metrics.Gauge(fmt.Sprintf("zymurgauge.%s.auxiliary_temperature,sensor_id=%s", c.Name,
			c.auxiliaryThermometer.GetID()), *c.Readings.AuxiliaryTemperature)
	}

	if c.Readings.ExternalTemperature != nil {
		c.metrics.Gauge(fmt.Sprintf("zymurgauge.%s.external_temperature,sensor_id=%s", c.Name,
			c.externalThermometer.GetID()), *c.Readings.ExternalTemperature)
	}

	if c.Readings.HydrometerGravity != nil {
		c.metrics.Gauge(fmt.Sprintf("zymurgauge.%s.hydrometer_gravity,sensor_id=%s", c.Name,
			c.hydrometer.GetID()), *c.Readings.HydrometerGravity)
	}
}

func (c *Chamber) sendToBrewFather(ctx context.Context) error {
	l := brewfather.LogEntry{
		DeviceName:      c.Name,
		Beer:            c.CurrentBatch.Recipe.Name,
		TemperatureUnit: "C",
		GravityUnit:     "G",
	}

	if c.Readings.BeerTemperature != nil {
		l.BeerTemperature = fmt.Sprintf("%f", *c.Readings.BeerTemperature)
	}

	if c.Readings.AuxiliaryTemperature != nil {
		l.AuxiliaryTemperature = fmt.Sprintf("%f", *c.Readings.AuxiliaryTemperature)
	}

	if c.Readings.ExternalTemperature != nil {
		l.ExternalTemperature = fmt.Sprintf("%f", *c.Readings.ExternalTemperature)
	}

	if c.Readings.HydrometerGravity != nil {
		l.Gravity = fmt.Sprintf("%f", *c.Readings.HydrometerGravity)
	}

	if err := c.service.Log(ctx, l); err != nil {
		return errors.Wrap(err, "could not log to Brewfather")
	}

	return nil
}
