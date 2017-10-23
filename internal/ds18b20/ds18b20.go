package ds18b20

import (
	"bufio"
	"errors"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/benjaminbartels/zymurgauge/internal/platform/safeclose"
)

const (
	prefix = "28-"
	slave  = "w1_slave"
)

// DevicePath is the location of the thermometer data on teh file system
var DevicePath = "/sys/bus/w1/devices/"

// GetThermometers returns a list of all Thermometers on the bus
func GetThermometers() ([]Thermometer, error) {
	var thermometers []Thermometer

	dir, err := os.Open(DevicePath)
	if err != nil {
		return nil, err
	}
	infos, err := dir.Readdir(-1)
	defer safeclose.Close(dir, &err)
	if err != nil {
		return nil, err
	}
	sort.Slice(infos, func(i, j int) bool { return infos[i].Name() < infos[j].Name() })
	for _, info := range infos {
		if strings.HasPrefix(info.Name(), prefix) {

			term := Thermometer{
				ID:   info.Name(),
				path: DevicePath,
			}

			thermometers = append(thermometers, term)
		}
	}
	return thermometers, nil
}

// GetThermometer returns a Thermometer by id
func GetThermometer(id string) (*Thermometer, error) {
	_, err := os.Stat(path.Join(DevicePath, id))
	if err != nil {
		return nil, err
	}

	return &Thermometer{
		ID:   id,
		path: DevicePath,
	}, nil
}

// Thermometer is a GPIO temperature probe
type Thermometer struct {
	ID   string
	path string
}

// ReadTemperature read the current temperature of the Thermometer
func (t *Thermometer) ReadTemperature() (*float64, error) {
	file, err := os.Open(path.Join(DevicePath, t.ID, slave))
	if err != nil {
		return nil, err
	}

	defer safeclose.Close(file, &err)

	r := bufio.NewReader(file)

	crcLine, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	crcLine = strings.TrimRight(crcLine, "\n")
	if !strings.HasSuffix(crcLine, "YES") {
		return nil, errors.New("CRC error")
	}

	dataLine, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	temp, err := strconv.ParseFloat(strings.Split(strings.TrimSpace(dataLine), "=")[1], 64)
	if err != nil {
		return nil, err
	}

	temp = temp / 1000

	return &temp, nil
}
