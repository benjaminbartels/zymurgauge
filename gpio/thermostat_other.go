// +build darwin dragonfly freebsd netbsd plan9 openbsd solaris windows

package gpio

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const mockData = "af 01 4b 46 7f ff 01 10 bc : crc=bc YES\naf 01 4b 46 7f ff 01 10 bc t=%d"

func (t *Thermostat) SetTemperature(temp *float64) error {
	t.target = temp
	if t.target != nil && !t.isPolling {
		go t.poll()
	} else if t.target == nil {
		t.quit <- true // ToDo: Check?
	}
	return nil
}

func (t *Thermostat) poll() {

	t.quit = make(chan bool)
	t.path = os.TempDir() + "zymurgauge/" + t.ThermometerID
	err := t.prepare()
	if err != nil {
		t.logger.Error(err)
		return
	}

	t.isPolling = true // ToDo: Make atomic?
	for {
		select {
		case <-t.quit:
			t.isPolling = false
			return
		case <-time.After(t.Interval):
			v, err := t.getTemperature()
			if err != nil {
				t.logger.Error(err)
			} else {
				if t.target != nil {

					if v > *t.target {
						t.logger.Debugf("Temperature above Target. Current: %f Target: %f", v, t.target)
						err = t.decreaseTemp()
					} else if v < *t.target {
						t.logger.Debugf("Temperature below Target. Current: %f Target: %f", v, t.target)
						err = t.increaseTemp()
					} else {
						t.logger.Debugf("Temperature equals Target. Current: %f Target: %f", v, t.target)
					}
					if err != nil {
						t.logger.Error(err)
					}

				}
			}
		}
	}
}

func (t *Thermostat) getTemperature() (float64, error) {

	data, err := ioutil.ReadFile(t.path + "/w1_slave")
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

func (t *Thermostat) increaseTemp() error {
	t.logger.Debug("Simulating heating")
	return t.updateTemp(true)
}

func (t *Thermostat) decreaseTemp() error {
	t.logger.Debug("Simulating cooling")
	return t.updateTemp(false)
}

func (t *Thermostat) updateTemp(up bool) error {

	data, err := ioutil.ReadFile(t.path + "/w1_slave")
	if err != nil {
		return err
	}

	temp, err := strconv.Atoi(strings.Split(strings.TrimSpace(string(data)), "=")[2])
	if err != nil {
		return err
	}

	if up {
		temp = temp + 100
	} else {
		temp = temp - 100
	}

	data = []byte(fmt.Sprintf(mockData, temp))
	err = os.MkdirAll(t.path, 0700)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(t.path+"/w1_slave", data, 0700)
	return err

}

func (t *Thermostat) prepare() error {
	data := []byte(fmt.Sprintf(mockData, 26937))
	err := os.MkdirAll(t.path, 0700)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(t.path+"/w1_slave", data, 0700)
	return err
}
