package chamber

import (
	"context"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/batch"
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
	ID                      string         `json:"id"`
	Name                    string         `json:"name"`
	DeviceConfigs           []DeviceConfig `json:"deviceConfigs"`
	ChillerKp               float64        `json:"chillerKp"`
	ChillerKi               float64        `json:"chillerKi"`
	ChillerKd               float64        `json:"chillerKd"`
	HeaterKp                float64        `json:"heaterKp"`
	HeaterKi                float64        `json:"heaterKi"`
	HeaterKd                float64        `json:"heaterKd"`
	CurrentBatch            *batch.Batch   `json:"currentBatch,omitempty"`
	CurrentFermentationStep int            `json:"currentFermentationStep"`
	ModTime                 time.Time      `json:"modTime"`
	logger                  *logrus.Logger
	thermometer             device.Thermometer
	hydrometer              device.Hydrometer
	chiller                 device.Actuator
	heater                  device.Actuator
	temperatureController   device.TemperatureController
	cancelFunc              context.CancelFunc
	runMutex                *sync.Mutex
}

type DeviceConfig struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Roles []string `json:"roles"`
}

// TODO: refactor to use generics in the future.

func (c *Chamber) Configure(configurator Configurator, logger *logrus.Logger) error {
	c.logger = logger

	var errs []error

	for _, deviceConfig := range c.DeviceConfigs {
		if err := c.configureDevice(configurator, deviceConfig); err != nil {
			errs = append(errs, err)
		}
	}

	c.temperatureController = pid.NewPIDTemperatureController(
		c.thermometer, c.chiller, c.heater, c.ChillerKp, c.ChillerKi, c.ChillerKd,
		c.HeaterKp, c.HeaterKi, c.HeaterKd, logger)

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

func (c *Chamber) StartFermentation(ctx context.Context, step int) error {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if c.CurrentBatch == nil {
		return ErrNoCurrentBatch
	}

	if step <= 0 || step > len(c.CurrentBatch.Fermentation.Steps) {
		return ErrInvalidStep
	}

	if c.cancelFunc != nil {
		c.cancelFunc()
	}

	temp := c.CurrentBatch.Fermentation.Steps[step-1].StepTemp
	ctx, cancelFunc := context.WithCancel(ctx)
	c.cancelFunc = cancelFunc

	go func() {
		if err := c.temperatureController.Run(ctx, temp); err != nil {
			c.logger.WithError(err).Errorf("could not run temperature controller for chamber %s", c.Name)
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
