package herms

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
)

type Configurator interface {
	CreateDs18b20(thermometerID string) (device.Thermometer, error)
	CreateGPIOActuator(pin string) (device.Actuator, error)
}
