package chamber

import (
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/device/raspberrypi"
	"github.com/pkg/errors"
)

func CreateThermometer(thermometerID string) (device.Thermometer, error) {
	return createPiThermometer(thermometerID)
}

func createPiThermometer(thermometerID string) (device.Thermometer, error) {
	ds18b20, err := raspberrypi.NewDs18b20(thermometerID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create new Ds18b20 thermometer %s", thermometerID)
	}

	return ds18b20, nil
}

func CreateActuator(pin string) (device.Actuator, error) {
	return createPiActuator(pin)
}

func createPiActuator(pin string) (device.Actuator, error) {
	actuator, err := raspberrypi.NewGPIOActuator(pin)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create new raspberry pi gpio actuator for pin %s", pin)
	}

	return actuator, nil
}
