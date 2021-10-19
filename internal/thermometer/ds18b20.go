package thermometer

import (
	"github.com/pkg/errors"
	"periph.io/x/periph/conn/onewire"
	"periph.io/x/periph/devices/ds18b20"
)

var _ Thermometer = (*Ds18b20)(nil)

type Ds18b20 struct {
	dev *ds18b20.Dev
}

func NewDs18b20(bus onewire.Bus, address onewire.Address, resolutionBits int) (*Ds18b20, error) {
	dev, err := ds18b20.New(bus, address, resolutionBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not create ds18b20 thermometer")
	}

	return &Ds18b20{
		dev: dev,
	}, nil
}

func (d *Ds18b20) GetTemperature() (float64, error) {
	temp, err := d.dev.LastTemp()
	if err != nil {
		return 0, errors.Wrap(err, "could not read ds18b20 thermometer")
	}

	return temp.Celsius(), nil
}
