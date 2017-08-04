// +build linux

package internal

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type TemperatureController struct {
	ThermometerID string        `json:"thermometerId"`
	Chiller       *Device       `json:"chiller"`
	Heater        *Device       `json:"heater"`
	Interval      time.Duration `json:"interval"`
	target        *float64
	quit          chan bool
	isPolling     bool
	state         state
}

const (
	path = "/sys/devices/w1_bus_master1/"
)

func (t *TemperatureController) SetTemperature(temp *float64) error {
	t.target = temp
	t.quit = make(chan bool)
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

func (t *TemperatureController) cool(on bool) {
	if t.Chiller != nil {

		err := rpio.Open()
		if err != nil {
			panic(err)
		}
		pin := rpio.Pin(*t.CoolerGPIO)
		if on {
			pin.High()
			fmt.Printf("Setting Chiller GPIO %d to High\n", t.Chiller.GPIO)
		} else {
			pin.Low()
			fmt.Printf("Setting Chiller GPIO %d to Low\n", t.Chiller.GPIO)
		}
		rpio.Close()
	} else {
		fmt.PrintLn("No Chiller Configured")
	}

}

func (t *TemperatureController) heat(on bool) {
	if t.Heater != nil {
		err := rpio.Open()
		if err != nil {
			panic(err)
		}
		pin := rpio.Pin(*t.HeaterGPIO)
		if on {
			pin.High()
			fmt.Printf("Setting Heater GPIO %d to High\n", t.Heater.GPIO)
		} else {
			pin.Low()
			fmt.Printf("Setting Heater GPIO %d to Low\n", t.Heater.GPIO)
		}
		rpio.Close()
	} else {
		fmt.PrintLn("No Heater Configured")
	}
}
