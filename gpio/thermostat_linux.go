// +build linux

package gpio

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/davecheney/gpio"
)

func (t *Thermostat) SetTemperature(temp *float64) err {
	t.target = temp

	if !t.isActive {

		t.isActive = true
		// ToDo: Implement
	}
}

func (t *Thermostat) getTemperature() (*float64, error) {
	data, err := ioutil.ReadFile(t.path + t.ThermometerID + "/w1_slave")
	if err != nil {
		return 0, err
	}

	temp, err := strconv.ParseFloat(strings.Split(strings.TrimSpace(string(data)), "=")[2], 64)
	if err != nil {
		return 0, err
	}

	temp = temp / 1000

	return temp, nil
}

func (t *Thermostat) cool(on bool) {

	gpio.ModePWM

	if on {
		pin, err := gpio.OpenPin(t.CoolerGPIO, gpio.ModeOutput)
		if err != nil {
			panic(err)
		}

		pin.Set()

		pin.Close()
	} else {
		pin, err := gpio.OpenPin(t.CoolerGPIO, gpio.ModeOutput)
		if err != nil {
			panic(err)
		}

		pin.Clear()

		pin.Close()
	}

}

func (t *Thermostat) heat(on bool) {

	if on {
		pin, err := gpio.OpenPin(t.HeaterGPIO, gpio.ModeOutput)
		if err != nil {
			panic(err)
		}

		pin.Set()

		pin.Close()
	} else {
		pin, err := gpio.OpenPin(t.HeaterGPIO, gpio.ModeOutput)
		if err != nil {
			panic(err)
		}

		pin.Clear()

		pin.Close()
	}

}
