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

	if t.path == "" {
		t.path = os.TempDir() + "zymurgauge/" + t.ThermometerID
		t.Logger.Debugf("Creating %s", t.path)
		data := []byte(fmt.Sprintf(mockData, 26937))
		err := os.MkdirAll(t.path, 0700)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(t.path+"/w1_slave", data, 0700)
		if err != nil {
			return err
		}

		go t.updateTemp()
	}

	t.target = temp
	if t.target != nil && !t.isPolling {
		go t.poll()
		go t.updateTemp()
	} else if t.target == nil && t.quit != nil {
		t.quit <- true // ToDo: Check?
	}
	return nil
}

func (t *Thermostat) poll() {

	t.quit = make(chan bool)

	t.isPolling = true // ToDo: Make atomic?
	for {
		select {
		case <-t.quit:
			t.isPolling = false
			return
		case <-time.After(t.Interval):
			v, err := t.getTemperature()
			if err != nil {
				t.Logger.Error(err)
			} else {
				if t.target != nil {
					if v > *t.target {
						t.Logger.Infof("Temperature above Target. Current: %f Target: %f", v, *t.target)
						t.state = COOLING
					} else if v < *t.target {
						t.Logger.Infof("Temperature below Target. Current: %f Target: %f", v, *t.target)
						t.state = HEATING
					} else {
						t.Logger.Infof("Temperature equals Target. Current: %f Target: %f", v, *t.target)
						t.state = OFF
					}
					if err != nil {
						t.Logger.Error(err)
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

func (t *Thermostat) updateTemp() {

	for {

		data, err := ioutil.ReadFile(t.path + "/w1_slave")
		if err != nil {
			t.Logger.Error(err)
		} else {

			if len(data) > 0 {

				temp, err := strconv.Atoi(strings.Split(strings.TrimSpace(string(data)), "=")[2])
				if err != nil {
					t.Logger.Error(err)
				} else {

					if t.state == HEATING {
						temp = temp + 100
					} else if t.state == COOLING {
						temp = temp - 100
					}

					data = []byte(fmt.Sprintf(mockData, temp))
					err = os.MkdirAll(t.path, 0700)
					if err != nil {
						t.Logger.Error(err)
					} else {
						err = ioutil.WriteFile(t.path+"/w1_slave", data, 0700)
						if err != nil {
							t.Logger.Error(err)
						}
					}
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}
