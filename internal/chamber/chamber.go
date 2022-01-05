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

// TODO: NEXT Send updates tp brew father

// Chamber represents an insulated box (fridge) with internal heating/cooling elements that reacts to changes in
// monitored temperatures, by correcting small deviations from your desired fermentation temperature.
type Chamber struct {
	ID                      string            `json:"id"`
	Name                    string            `json:"name"`
	DeviceConfigs           []DeviceConfig    `json:"deviceConfigs"`
	ChillerKp               float64           `json:"chillerKp"`
	ChillerKi               float64           `json:"chillerKi"`
	ChillerKd               float64           `json:"chillerKd"`
	HeaterKp                float64           `json:"heaterKp"`
	HeaterKi                float64           `json:"heaterKi"`
	HeaterKd                float64           `json:"heaterKd"`
	CurrentBatch            *brewfather.Batch `json:"currentBatch,omitempty"`
	CurrentFermentationStep int               `json:"currentFermentationStep"` // TODO: Is this used?
	ModTime                 time.Time         `json:"modTime"`
	logger                  *logrus.Logger
	thermometer             device.Thermometer
	hydrometer              device.Hydrometer
	chiller                 device.Actuator
	heater                  device.Actuator
	temperatureController   device.TemperatureController
	service                 brewfather.Service
	statusCh                chan pid.Status
	tiltColor               string
	cancelFunc              context.CancelFunc
	runMutex                *sync.Mutex
}

type DeviceConfig struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Roles []string `json:"roles"`
}

// TODO: refactor to use generics in the future.

func (c *Chamber) Configure(configurator Configurator, service brewfather.Service, logger *logrus.Logger) error {
	c.logger = logger

	var errs []error

	for _, deviceConfig := range c.DeviceConfigs {
		if err := c.configureDevice(configurator, deviceConfig); err != nil {
			errs = append(errs, err)
		}
	}

	c.service = service

	c.statusCh = make(chan pid.Status, 1)

	c.temperatureController = pid.NewPIDTemperatureController(
		c.thermometer, c.chiller, c.heater, c.ChillerKp, c.ChillerKi, c.ChillerKd,
		c.HeaterKp, c.HeaterKi, c.HeaterKd, logger, pid.SetStatusChannel(c.statusCh))

	c.runMutex = &sync.Mutex{}

	if errs != nil {
		return &InvalidConfigurationError{configErrors: errs}
	}

	return nil
}

func (c *Chamber) configureDevice(configurator Configurator, deviceConfig DeviceConfig) error {
	var (
		createdDevice interface{}
		err           error
	)

	switch deviceConfig.Type {
	case "ds18b20":
		if createdDevice, err = configurator.CreateDs18b20(deviceConfig.ID); err != nil {
			return errors.Wrapf(err, "could not create new Ds18b20 %s", deviceConfig.ID)
		}
	case "tilt":
		if createdDevice, err = configurator.CreateTilt(tilt.Color(deviceConfig.ID)); err != nil {
			return errors.Wrapf(err, "could not create new %s Tilt", deviceConfig.ID)
		}
		c.tiltColor = deviceConfig.ID
	case "gpio":
		if createdDevice, err = configurator.CreateGPIOActuator(deviceConfig.ID); err != nil {
			return errors.Wrapf(err, "could not create new GPIO %s", deviceConfig.ID)
		}
	default:
		return errors.Errorf("invalid device type '%s'", deviceConfig.Type)
	}

	if err := c.assignDevice(createdDevice, deviceConfig.Roles); err != nil {
		return err
	}

	return nil
}

// TODO: NEXT make types match brewfather

func (c *Chamber) assignDevice(d interface{}, roles []string) error {
	// type assertions will not fail
	for _, role := range roles {
		switch role {
		case "thermometer":
			c.thermometer, _ = d.(device.Thermometer)
		case "hydrometer":
			c.hydrometer, _ = d.(device.Hydrometer)
		case "chiller":
			c.chiller, _ = d.(device.Actuator)
		case "heater":
			c.heater, _ = d.(device.Actuator)
		default:
			return errors.Errorf("invalid device role '%s'", role)
		}
	}

	return nil
}

func (c *Chamber) StartFermentation(ctx context.Context, stepId string) error {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if c.CurrentBatch == nil {
		return ErrNoCurrentBatch
	}

	var step *brewfather.FermentationStep
	for _, s := range c.CurrentBatch.Fermentation.Steps {
		if s.Type == stepId {
			step = &s
			break
		}
	}

	if step == nil {
		return ErrInvalidStep
	}

	if c.cancelFunc != nil {
		c.cancelFunc()
	}

	temp := step.StepTemp
	ctx, cancelFunc := context.WithCancel(ctx)
	c.cancelFunc = cancelFunc

	if c.tiltColor != "" { // TODO: NEXT dont refernce tilt
		go func() {
			c.sendData(ctx)
			for {
				timer := time.NewTimer(15 * time.Minute)
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

	go func() {
		if err := c.temperatureController.Run(ctx, temp); err != nil {
			c.logger.WithError(err).Errorf("could not run temperature controller for chamber %s", c.Name)
		}
	}()

	return nil
}

func (c *Chamber) sendData(ctx context.Context) {
	temperature, err := c.thermometer.GetTemperature()
	if err != nil {
		c.logger.WithError(err).Error("could not read temperature")
	}

	specificGravity, err := c.hydrometer.GetSpecificGravity()
	if err != nil {
		c.logger.WithError(err).Error("could not get specific gravity")
	}

	log := brewfather.LogEntry{
		DeviceName:      c.Name,
		BeerTemperature: fmt.Sprintf("%f", temperature),
		TemperatureUnit: "C",
		Gravity:         fmt.Sprintf("%f", specificGravity),
		GravityUnit:     "G",
		Beer:            c.CurrentBatch.Name,
	}

	c.logger.Debugf("Sending Data to Brewfather: Beer: %s Temperature: %s SpecificGravity: %s",
		log.Beer, log.BeerTemperature, log.Gravity)
	if err := c.service.Log(ctx, log); err != nil {
		c.logger.WithError(err).Error("could log tilt data")
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
