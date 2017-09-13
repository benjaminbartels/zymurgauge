package internal

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func (t *TemperatureController) Equals(o *TemperatureController) bool {
	if t.ThermometerID != o.ThermometerID {
		return false
	} else if t.Chiller.Equals(o.Chiller) {
		return false
	} else if t.Heater.Equals(o.Heater) {
		return false
	} else if t.Interval != o.Interval {
		return false
	}
	return true
}

func (t *TemperatureController) poll() {
	t.isPolling = true // ToDo: Make atomic?

	t.quit = make(chan bool)

	for {
		select {
		case <-t.quit:
			t.isPolling = false
			return
		case <-time.After(t.Interval):
			if t.target != nil {
				v, err := t.getTemperature()
				if err != nil {
					fmt.Println(err)
					continue
				}

				t.evaluateTemperature(v)
			}
		}
	}
}

func (t *TemperatureController) getTemperature() (float64, error) {
	data, err := ioutil.ReadFile(path + t.ThermometerID + "/w1_slave")
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

func (t *TemperatureController) evaluateTemperature(temperature float64) {

	if temperature > *t.target {
		fmt.Printf("Temperature above Target. Current: %f Target: %f\n", temperature, *t.target)

		if t.state != COOLING {
			t.heat(false)
			if t.cool(true) {
				t.state = COOLING
			}
		}

	} else if temperature < *t.target {
		fmt.Printf("Temperature below Target. Current: %f Target: %f\n", temperature, *t.target)

		if t.state != HEATING {
			t.cool(false)
			t.heat(true)
			t.state = HEATING
		}

	} else {
		fmt.Printf("Temperature equals Target. Current: %f Target: %f\n", temperature, *t.target)

		if t.state != OFF {
			t.cool(false)
			t.heat(false)
			t.state = OFF
		}
	}
}

type state int

const (
	OFF state = 1 + iota
	COOLING
	HEATING
)
