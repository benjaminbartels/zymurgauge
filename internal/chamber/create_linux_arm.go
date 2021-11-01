package chamber

import (
	"github.com/benjaminbartels/zymurgauge/internal/device/pid"
	"github.com/benjaminbartels/zymurgauge/internal/device/raspberrypi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func CreateThermometer(thermometerID string) (device.Thermometer, error) {
	return createPiThermometer(thermometerID)
}

func createPiThermometer(thermometerID string) (device.Thermometer, error) {
	thermometer, err := raspberrypi.NewDs18b20(thermometerID)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new Ds18b20 thermometer")
	}

	ds18b20, err := raspberrypi.NewDs18b20(c.ThermometerID)
	if err != nil {
		return errors.Wrapf(err, "could not create new Ds18b20 thermometer %s", c.ThermometerID)
	}

	return ds18b20, nil
}

func CreateActuator(pin string) (device.Actuator, error) {
	return createPiActuator(pin)
}

func createPiActuator(pin string) (device.Actuator, error) {
	actuator, err := raspberrypi.NewGPIOActuator(pin)
	if err != nil {
		return errors.Wrapf(err, "could not create new raspberry pi gpio actuator for pin %s", pin)
	}

	return actuator, nil
}
