// +build linux

package gpio

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"fmt"

	"github.com/stianeikeland/go-rpio"
)

func (t *Thermostat) SetTemperature(temp *float64) error {
	t.target = temp
	if t.target != nil && !t.isPolling {
		go t.poll()
	} else if t.target == nil {
		t.quit <- true // ToDo: Check?
	}
	return nil
}

func init() {
	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	rpio.Pin(0).Low()
	rpio.Pin(4).Low()
	rpio.Pin(7).Low()
	rpio.Pin(8).Low()
	rpio.Pin(9).Low()
	rpio.Pin(10).Low()
	rpio.Pin(11).Low()
	rpio.Pin(14).Low()
	rpio.Pin(15).Low()
	rpio.Pin(17).Low()
	rpio.Pin(18).Low()
	rpio.Pin(21).Low()
	rpio.Pin(22).Low()
	rpio.Pin(23).Low()
	rpio.Pin(24).Low()
	rpio.Pin(25).Low()
	rpio.Close()
}

func (t *Thermostat) poll() {

	t.quit = make(chan bool)
	t.path = "/sys/devices/w1_bus_master1/"

	t.isPolling = true // ToDo: Make atomic?

	for {
		select {
		case <-t.quit:
			t.isPolling = false
			return
		case <-time.After(t.Interval):
			if t.target != nil {
				v, err := t.getTemperature()
				if err != nil {
					t.Logger.Error(err)
					continue
				}

				t.evaluateTemperature(v)
			}
		}
	}
}

func (t *Thermostat) getTemperature() (float64, error) {
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

func (t *Thermostat) evaluateTemperature(temperature float64) {

	if temperature > *t.target {
		t.Logger.Infof("Temperature above Target. Current: %f Target: %f", temperature, *t.target)

		if t.state != COOLING {
			t.cool(true)
			t.heat(false)
			t.state = COOLING
		}

	} else if temperature < *t.target {
		t.Logger.Infof("Temperature below Target. Current: %f Target: %f", temperature, *t.target)

		if t.state != HEATING {
			t.cool(false)
			t.heat(true)
			t.state = HEATING
		}

	} else {
		t.Logger.Infof("Temperature equals Target. Current: %f Target: %f", temperature, *t.target)

		if t.state != OFF {
			t.cool(false)
			t.heat(false)
			t.state = OFF
		}
	}
}

func (t *Thermostat) cool(on bool) {

	fmt.Println("COOL = ", on)

	if t.CoolerGPIO != nil {

		err := rpio.Open()
		if err != nil {
			panic(err)
		}
		pin := rpio.Pin(*t.CoolerGPIO)
		if on {
			pin.High()
		} else {
			pin.Low()
		}
		rpio.Close()
	}

}

func (t *Thermostat) heat(on bool) {
	if t.HeaterGPIO != nil {

		err := rpio.Open()
		if err != nil {
			panic(err)
		}
		pin := rpio.Pin(*t.HeaterGPIO)
		if on {
			pin.High()
		} else {
			pin.Low()
		}
		rpio.Close()
	}
}
