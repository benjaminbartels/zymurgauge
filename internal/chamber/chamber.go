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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const brewfatherLogInterval = 15 * time.Minute

// Chamber represents an insulated box (fridge) with internal heating/cooling elements that reacts to changes in
// monitored temperatures, by correcting small deviations from your desired fermentation temperature.
type Chamber struct {
	ID                      string                  `json:"id"` // TODO: omit empty?
	Name                    string                  `json:"name"`
	DeviceConfig            DeviceConfig            `json:"deviceConfigs"` // TODO: make is a struct and not a list
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
	beerThermometer         device.Thermometer
	auxiliaryThermometer    device.Thermometer
	externalThermometer     device.Thermometer
	hydrometer              device.Hydrometer
	chiller                 device.Actuator
	heater                  device.Actuator
	temperatureController   device.TemperatureController
	service                 brewfather.Service
	logToBrewfather         bool
	cancelFunc              context.CancelFunc
	runMutex                *sync.Mutex
}

type DeviceConfig struct {
	ChillerGPIO               string `json:"chillerGpio"`
	HeaterGPIO                string `json:"heaterGpio"`
	BeerThermometerType       string `json:"beerThermometerType"`
	BeerThermometerID         string `json:"beerThermometerId"`
	AuxiliaryThermometerType  string `json:"auxiliaryThermometerType"`
	AuxiliaryThermometerID    string `json:"auxiliaryThermometerId"`
	ExternalThermometerType   string `json:"externalThermometerType"`
	ExternalThermometerID     string `json:"externalThermometerId"`
	HydrometerThermometerType string `json:"hydrometerThermometerType"`
	HydrometerThermometerID   string `json:"hydrometerThermometerId"`
}

type Readings struct {
	BeerTemperature      float64 `json:"beerTemperature"`
	AuxiliaryTemperature float64 `json:"auxiliaryTemperature"`
	ExternalTemperature  float64 `json:"externalTemperature"`
	HydrometerGravity    float64 `json:"hydrometerGravity"`
}

// TODO: refactor to use generics in the future.

func (c *Chamber) Configure(configurator Configurator, service brewfather.Service, logToBrewfather bool,
	logger *logrus.Logger) error {
	c.logger = logger

	c.service = service
	c.logToBrewfather = logToBrewfather

	errs := c.configureDevices(configurator, c.DeviceConfig)

	c.temperatureController = pid.NewPIDTemperatureController(
		c.beerThermometer, c.chiller, c.heater, c.ChillerKp, c.ChillerKi, c.ChillerKd,
		c.HeaterKp, c.HeaterKi, c.HeaterKd, logger)

	c.runMutex = &sync.Mutex{}

	if errs != nil {
		return &InvalidConfigurationError{configErrors: errs}
	}

	return nil
}

func (c *Chamber) configureDevices(configurator Configurator, config DeviceConfig) []error {
	var errs []error

	errs = append(errs, c.configureActuators(configurator, config)...)

	errs = append(errs, c.configureThermometers(configurator, config)...)

	if config.HydrometerThermometerType != "" {
		h, err := getHydrometer(configurator, config.HydrometerThermometerType,
			config.HydrometerThermometerID)
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

	if c.logToBrewfather {
		go func() {
			c.sendData(ctx)

			for {
				timer := time.NewTimer(brewfatherLogInterval)
				defer timer.Stop()

				select {
				case <-timer.C:
					c.sendData(ctx)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

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

func (c *Chamber) getStep(stepID string) *brewfather.FermentationStep {
	var step *brewfather.FermentationStep

	for i := range c.CurrentBatch.Fermentation.Steps {
		if c.CurrentBatch.Fermentation.Steps[i].Type == stepID {
			step = &c.CurrentBatch.Fermentation.Steps[i]

			break
		}
	}

	return step
}

func (c *Chamber) sendData(ctx context.Context) {
	log := brewfather.LogEntry{
		DeviceName:      c.Name,
		Beer:            c.CurrentBatch.Name,
		TemperatureUnit: "C",
		GravityUnit:     "G",
	}

	if t, err := c.getBeerTemperature(); err != nil {
		c.logger.WithError(err).Error("could not get beer temperature")
	} else {
		log.BeerTemperature = fmt.Sprintf("%f", t)
	}

	if t, err := c.getAuxiliaryTemperature(); err != nil {
		c.logger.WithError(err).Error("could not get auxiliary temperature")
	} else {
		log.AuxiliaryTemperature = fmt.Sprintf("%f", t)
	}

	if t, err := c.getExternalTemperature(); err != nil {
		c.logger.WithError(err).Error("could not get external temperature")
	} else {
		log.ExternalTemperature = fmt.Sprintf("%f", t)
	}

	if t, err := c.getHydrometerGravity(); err != nil {
		c.logger.WithError(err).Error("could not get external temperature")
	} else {
		log.Gravity = fmt.Sprintf("%f", t)
	}

	c.logger.Debugf("Sending Data to Brewfather: Beer: %s Temperature: %s Gravity: %s",
		log.Beer, log.BeerTemperature, log.Gravity)

	if err := c.service.Log(ctx, log); err != nil {
		c.logger.WithError(err).Error("could not log tilt data")
	}
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

func (c *Chamber) UpdateReadings() {
	c.Readings = &Readings{}

	if t, err := c.getBeerTemperature(); err != nil {
		c.logger.WithError(err).Error("could not get reading for beer temperature")
	} else {
		c.Readings.BeerTemperature = t
	}

	if t, err := c.getAuxiliaryTemperature(); err != nil {
		c.logger.WithError(err).Error("could not get reading for auxiliary temperature")
	} else {
		c.Readings.AuxiliaryTemperature = t
	}

	if t, err := c.getExternalTemperature(); err != nil {
		c.logger.WithError(err).Error("could not get reading for external temperature")
	} else {
		c.Readings.ExternalTemperature = t
	}

	if t, err := c.getHydrometerGravity(); err != nil {
		c.logger.WithError(err).Error("could not get reading for hydrometer gravity")
	} else {
		c.Readings.HydrometerGravity = t
	}
}

func (c *Chamber) getBeerTemperature() (float64, error) {
	if c.beerThermometer == nil {
		return 0, errors.New("beer thermometer is nil")
	}

	t, err := c.beerThermometer.GetTemperature()
	if err != nil {
		return 0, errors.Wrap(err, "could not get beer temperature")
	}

	return t, nil
}

func (c *Chamber) getAuxiliaryTemperature() (float64, error) {
	if c.auxiliaryThermometer == nil {
		return 0, errors.New("auxiliary thermometer is nil")
	}

	t, err := c.auxiliaryThermometer.GetTemperature()
	if err != nil {
		return 0, errors.Wrap(err, "could not get auxiliary temperature")
	}

	return t, nil
}

func (c *Chamber) getExternalTemperature() (float64, error) {
	if c.externalThermometer == nil {
		return 0, errors.New("external thermometer is nil")
	}

	t, err := c.externalThermometer.GetTemperature()
	if err != nil {
		return 0, errors.Wrap(err, "could not get external temperature")
	}

	return t, nil
}

func (c *Chamber) getHydrometerGravity() (float64, error) {
	if c.hydrometer == nil {
		return 0, errors.New("hydrometer is nil")
	}

	t, err := c.hydrometer.GetGravity()
	if err != nil {
		return 0, errors.Wrap(err, "could not get hydrometer gravity")
	}

	return t, nil
}
