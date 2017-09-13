// +build darwin dragonfly freebsd netbsd plan9 openbsd solaris windows

package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
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
	isUpdating    bool
}

const (
	path     = "./"
	filename = "w1_slave"
	mockData = "af 01 4b 46 7f ff 01 10 bc : crc=bc YES\naf 01 4b 46 7f ff 01 10 bc t=%d"
)

func (t *TemperatureController) SetTemperature(temp *float64) error {

	// detect if file exists
	var _, err = os.Stat(path + t.ThermometerID + "/" + filename)

	// create file if not exists
	if os.IsNotExist(err) {

		fmt.Printf("Creating %s", path+t.ThermometerID+"/"+filename+"\n")
		data := []byte(fmt.Sprintf(mockData, 26937))
		err := os.MkdirAll(path+t.ThermometerID, 0700)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path+t.ThermometerID+"/"+filename, data, 0700)
		if err != nil {
			return err
		}
	}

	if !t.isUpdating {
		go t.updateTemp()
	}

	t.target = temp
	t.quit = make(chan bool)
	if t.target != nil && !t.isPolling {
		go t.poll()
	} else if t.target == nil {
		t.quit <- true // ToDo: Check?
	}
	return nil
}

func (t *TemperatureController) cool(on bool) bool {
	if t.Chiller != nil {
		if on {
			fmt.Printf("Setting Chiller GPIO %d to High\n", t.Chiller.GPIO)
		} else {
			fmt.Printf("Setting Chiller GPIO %d to Low\n", t.Chiller.GPIO)
		}
	} else {
		fmt.Println("No Chiller Configured")
	}
	// fis this
	return true
}

func (t *TemperatureController) heat(on bool) {
	if t.Heater != nil {
		if on {
			fmt.Printf("Setting Heater GPIO %d to High\n", t.Heater.GPIO)
		} else {
			fmt.Printf("Setting Heater GPIO %d to Low\n", t.Heater.GPIO)
		}
	} else {
		fmt.Println("No Heater Configured")
	}
}

func (t *TemperatureController) updateTemp() {
	t.isUpdating = true
	for {
		data, err := ioutil.ReadFile(path + t.ThermometerID + "/" + filename)
		if err != nil {
			fmt.Println(err)
		} else {

			if len(data) > 0 {

				temp, err := strconv.Atoi(strings.Split(strings.TrimSpace(string(data)), "=")[2])
				if err != nil {
					fmt.Println(err)
				} else {

					if t.state == HEATING {
						temp = temp + 100
					} else if t.state == COOLING {
						temp = temp - 100
					}

					data = []byte(fmt.Sprintf(mockData, temp))
					err = ioutil.WriteFile(path+t.ThermometerID+"/"+filename, data, 0700)
					if err != nil {
						fmt.Println(err)
					}

				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}
