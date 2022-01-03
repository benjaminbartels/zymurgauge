package chamber

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
)

type Controller interface {
	Repo
	StartFermentation(chamberID string, step int) error
	StopFermentation(chamberID string) error
}

type Configurator interface {
	CreateDs18b20(thermometerID string) (device.Thermometer, error)
	CreateTilt(color tilt.Color) (device.ThermometerAndHydrometer, error)
	CreateGPIOActuator(pin string) (device.Actuator, error)
}

type Repo interface {
	GetAll() ([]*Chamber, error) // TODO: add ctx?
	Get(id string) (*Chamber, error)
	Save(c *Chamber) error
	Delete(id string) error
}
