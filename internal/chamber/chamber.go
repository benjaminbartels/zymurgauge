package chamber

import (
	"context"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/batch"
	"github.com/benjaminbartels/zymurgauge/internal/configurator"
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/pid"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidDeviceType = errors.New("invalid device type")
	ErrInvalidDeviceRole = errors.New("invalid device role")
	ErrInvalidStep       = errors.New("invalid step")
	ErrNoCurrentBatch    = errors.New("chamber does not have a current batch")
	ErrNotFermenting     = errors.New("fermentation has not started")
)

type Fermentor interface {
	StartFermentation(ctx context.Context, step int) error
	StopFermentation() error
}

type DeviceConfig struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Roles []string `json:"roles"`
}

// Chamber represents an insulated box (fridge) with internal heating/cooling elements that reacts to changes in
// monitored temperatures, by correcting small deviations from your desired fermentation temperature.
type Chamber struct {
	ID                      string         `json:"id"`
	Name                    string         `json:"name"`
	DeviceConfigs           []DeviceConfig `json:"thermometer"`
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

// TODO: refactor to use generics in the future.

func (c *Chamber) Configure(configurator configurator.ConfiguratorIface, logger *logrus.Logger) error {
	c.logger = logger

	var (
		createdDevice interface{}
		err           error
	)

	for _, deviceConfig := range c.DeviceConfigs {
		switch deviceConfig.Type {
		case "ds18b20":
			createdDevice, err = configurator.CreateDs18b20(deviceConfig.ID)
			if err != nil {
				return errors.Wrapf(err, "could not create new Ds18b20 thermometer %s", c.ID)
			}

		case "tilt":
			createdDevice, err = configurator.CreateTilt(tilt.Color(deviceConfig.ID))
			if err != nil {
				return errors.Wrapf(err, "could not create new Tilt %s", c.ID)
			}

		case "gpio":
			createdDevice, err = configurator.CreateGPIOActuator(deviceConfig.ID)
			if err != nil {
				return errors.Wrapf(err, "could not create new Tilt %s", c.ID)
			}

		default:
			return ErrInvalidDeviceType
		}

		if err := c.assign(createdDevice, deviceConfig.Roles); err != nil {
			return err
		}
	}

	// chiller, err := CreateGPIOActuator(c.ChillerPin)
	// if err != nil {
	// 	return errors.Wrapf(err, "could not create new chiller gpio actuator for pin %s", c.ChillerPin)
	// }

	// heater, err := CreateGPIOActuator(c.HeaterPin)
	// if err != nil {
	// 	return errors.Wrapf(err, "could not create new heater gpio actuator for pin %s", c.HeaterPin)
	// }

	c.temperatureController = pid.NewPIDTemperatureController(
		c.thermometer, c.chiller, c.heater, c.ChillerKp, c.ChillerKi, c.ChillerKd,
		c.HeaterKp, c.HeaterKi, c.HeaterKd, logger)

	c.runMutex = &sync.Mutex{}

	return nil
}

func (c *Chamber) assign(d interface{}, roles []string) error {
	for _, role := range roles {
		var ok bool

		switch role {
		case "thermometer":
			if c.thermometer, ok = d.(device.Thermometer); !ok {
				return errors.New("device is not a thermometer")
			}
		case "hydrometer":
			if c.hydrometer, ok = d.(device.Hydrometer); !ok {
				return errors.New("device is not a hydrometer")
			}
		case "chiller":
			if c.chiller, ok = d.(device.Actuator); !ok {
				return errors.New("device is not a hydrometer")
			}
		case "heater":
			if c.heater, ok = d.(device.Actuator); !ok {
				return errors.New("device is not a hydrometer")
			}
		default:
			return ErrInvalidDeviceRole
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

	return nil
}
