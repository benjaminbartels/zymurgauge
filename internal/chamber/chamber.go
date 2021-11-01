package chamber

import (
	"context"
	"sync"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/batch"
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/pid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidStep   = errors.New("invalid step")
	ErrNotFermenting = errors.New("fermentation has not started")
)

// Chamber represents an insulated box (fridge) with internal heating/cooling elements that reacts to changes in
// monitored temperatures, by correcting small deviations from your desired fermentation temperature.
type Chamber struct {
	ID                      string       `json:"id"`
	Name                    string       `json:"name"`
	ThermometerID           string       `json:"thermometerId"`
	ChillerPin              string       `json:"chillerPin"`
	HeaterPin               string       `json:"heaterPin"`
	ChillerKp               float64      `json:"chillerKp"`
	ChillerKi               float64      `json:"chillerKi"`
	ChillerKd               float64      `json:"chillerKd"`
	HeaterKp                float64      `json:"heaterKp"`
	HeaterKi                float64      `json:"heaterKi"`
	HeaterKd                float64      `json:"heaterKd"`
	CurrentBatch            *batch.Batch `json:"currentBatch,omitempty"`
	CurrentFermentationStep int          `json:"currentFermentationStep"`
	ModTime                 time.Time    `json:"modTime"`
	mainCtx                 context.Context
	logger                  *logrus.Logger
	temperatureController   device.TemperatureController
	cancelFunc              context.CancelFunc
	runMutex                *sync.Mutex
}

func (c *Chamber) Configure(ctx context.Context, logger *logrus.Logger) error {
	c.mainCtx = ctx
	c.logger = logger

	createThermometerFunc := CreateThermometer

	thermometer, err := createThermometerFunc(c.ThermometerID)
	if err != nil {
		return errors.Wrapf(err, "could not create new thermometer %s", c.ThermometerID)
	}

	createActuatorFunc := CreateActuator

	chiller, err := createActuatorFunc(c.ChillerPin)
	if err != nil {
		return errors.Wrapf(err, "could not create new chiller gpio actuator for pin %s", c.ChillerPin)
	}

	heater, err := createActuatorFunc(c.HeaterPin)
	if err != nil {
		return errors.Wrapf(err, "could not create new heater gpio actuator for pin %s", c.HeaterPin)
	}

	c.temperatureController = pid.NewPIDTemperatureController(
		thermometer, chiller, heater, c.ChillerKp, c.ChillerKi, c.ChillerKd,
		c.HeaterKp, c.HeaterKi, c.HeaterKd, logger)

	c.runMutex = &sync.Mutex{}

	return nil
}

func (c *Chamber) StartFermentation(step int) error {
	c.runMutex.Lock()
	defer c.runMutex.Unlock()

	if step < 0 || step >= len(c.CurrentBatch.Fermentation.Steps) {
		return ErrInvalidStep
	}

	if c.cancelFunc != nil {
		c.cancelFunc()
	}

	temp := c.CurrentBatch.Fermentation.Steps[step].StepTemp
	ctx, cancelFunc := context.WithCancel(c.mainCtx)
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
