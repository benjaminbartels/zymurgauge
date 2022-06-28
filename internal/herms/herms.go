package herms

import (
	"context"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/batch"
	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/platform/metrics"
	"github.com/benjaminbartels/zymurgauge/internal/temperaturecontrol/pid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type HERMS struct {
	ID                                 string        `json:"id,omitempty"`
	Name                               string        `json:"name"`
	DeviceConfig                       DeviceConfig  `json:"deviceConfig"`
	CurrentBatch                       *batch.Detail `json:"currentBatch,omitempty"`
	ModTime                            time.Time     `json:"modTime"`
	Readings                           *Readings     `json:"readings,omitempty"`
	logger                             *logrus.Logger
	metrics                            metrics.Metrics
	hotLiquorTankThermometer           device.Thermometer
	brewKettleThermometer              device.Thermometer
	mashThermometer                    device.Thermometer
	hotLiquorTankHeater                device.Actuator
	brewKettleHeater                   device.Actuator
	hotLiquorTankTemperatureController device.TemperatureController
	brewKettleTemperatureController    device.TemperatureController
	service                            brewfather.Service
	cancelFunc                         context.CancelFunc
	readingsUpdateInterval             time.Duration
	runMutex                           *sync.RWMutex
	readingsMutex                      *sync.Mutex
}

type DeviceConfig struct {
	HotLiquorTankHeaterKp      float64 `json:"hotLiquorTankHeaterKp"`
	HotLiquorTankHeaterKi      float64 `json:"hotLiquorTankHeaterKi"`
	HotLiquorTankHeaterKd      float64 `json:"hotLiquorTankHeaterKd"`
	BrewKettleHeaterKp         float64 `json:"brewKettleHeaterKp"`
	BrewKettleHeaterKi         float64 `json:"brewKettleHeaterKi"`
	BrewKettleHeaterKd         float64 `json:"brewKettleHeaterKd"`
	HotLiquorTankHeaterGPIO    string  `json:"hotLiquorTankHeaterGpio"`
	BrewKettleHeaterGPIO       string  `json:"brewKettleHeaterGpio"`
	HotLiquorTankThermometerID string  `json:"hotLiquorTankThermometerId"`
	BrewKettleThermometerID    string  `json:"brewKettleThermometerId"`
	MashThermometerID          string  `json:"mashThermometerId"`
}

type Readings struct {
	HotLiquorTankTemperature *float64 `json:"hotLiquorTankTemperature,omitempty"`
	BrewKettleTemperature    *float64 `json:"brewKettleTemperature,omitempty"`
	MashTemperature          *float64 `json:"mashTemperature,omitempty"`
}

func (h *HERMS) Configure(configurator Configurator, service brewfather.Service,
	logger *logrus.Logger, metrics metrics.Metrics, readingsUpdateInterval time.Duration,
) error {
	h.service = service
	h.logger = logger
	h.metrics = metrics
	h.readingsUpdateInterval = readingsUpdateInterval

	errs := h.configureDevices(configurator, h.DeviceConfig)

	h.hotLiquorTankTemperatureController = pid.NewController(h.mashThermometer, h.hotLiquorTankHeater,
		h.DeviceConfig.HotLiquorTankHeaterKp, h.DeviceConfig.HotLiquorTankHeaterKi, h.DeviceConfig.HotLiquorTankHeaterKd,
		logger)

	h.runMutex = &sync.RWMutex{}

	if errs != nil {
		return &InvalidConfigurationError{configErrors: errs}
	}

	return nil
}

func (h *HERMS) configureDevices(configurator Configurator, config DeviceConfig) []error {
	var errs []error

	errs = append(errs, h.configureActuators(configurator, config)...)

	errs = append(errs, h.configureThermometers(configurator, config)...)

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func (h *HERMS) configureActuators(configurator Configurator, config DeviceConfig) []error {
	var errs []error

	var err error

	h.hotLiquorTankHeater, err = configurator.CreateGPIOActuator(config.HotLiquorTankHeaterGPIO)
	if err != nil {
		errs = append(errs, errors.Wrapf(err, "could not create new GPIO %s for hot liquor tank heater ",
			config.HotLiquorTankHeaterGPIO))
	}

	h.brewKettleHeater, err = configurator.CreateGPIOActuator(config.BrewKettleHeaterGPIO)
	if err != nil {
		errs = append(errs, errors.Wrapf(err, "could not create new GPIO %s for brew kettle heater",
			config.BrewKettleHeaterGPIO))
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func (h *HERMS) configureThermometers(configurator Configurator, config DeviceConfig) []error {
	var errs []error

	var err error

	h.mashThermometer, err = configurator.CreateDs18b20(config.MashThermometerID)
	if err != nil {
		errs = append(errs, errors.Wrap(err, "could not configure mash thermometer"))
	}

	h.brewKettleThermometer, err = configurator.CreateDs18b20(config.BrewKettleThermometerID)
	if err != nil {
		errs = append(errs, errors.Wrap(err, "could not configure brew kettle thermometer"))
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
