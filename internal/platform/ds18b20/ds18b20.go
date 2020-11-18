package ds18b20

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	slave = "w1_slave"
	// defaultDevicePath is the location of the thermometer data on the file system.
	defaultDevicePath = "/sys/bus/w1/devices/"
)

// NewThermometer Create a new ds18b20 Thermometer.
func NewThermometer(id string) (*Thermometer, error) {
	return NewWithDevicePath(id, defaultDevicePath)
}

// NewWithDevicePath creates a new Thermometer using the given devicePath
// Usually used for testing.
func NewWithDevicePath(id, devicePath string) (*Thermometer, error) {
	_, err := os.Stat(path.Join(devicePath, id))
	if err != nil {
		return nil, errors.Wrap(err, "could not file information")
	}

	return &Thermometer{
		ID:   id,
		path: devicePath,
	}, nil
}

// Thermometer is a GPIO temperature probe.
type Thermometer struct {
	ID   string
	path string
}

// ReadTemperature read the current temperature of the Thermometer.
func (t *Thermometer) Read() (float64, error) {
	file, err := os.Open(path.Join(t.path, t.ID, slave))
	if err != nil {
		return 0, errors.Wrap(err, "could not open file")
	}

	defer file.Close()

	r := bufio.NewReader(file)

	crcLine, err := r.ReadString('\n')
	if err != nil {
		return 0, errors.Wrap(err, "could not read crc")
	}

	crcLine = strings.TrimRight(crcLine, "\n")
	if !strings.HasSuffix(crcLine, "YES") {
		return 0, errors.Wrap(err, "crc is invalid")
	}

	dataLine, err := r.ReadString('\n')
	if err != nil {
		return 0, errors.Wrap(err, "could not read data")
	}

	temp, err := strconv.ParseFloat(strings.Split(strings.TrimSpace(dataLine), "=")[1], 64)
	if err != nil {
		return 0, errors.Wrap(err, "could not parse temperature value")
	}

	temp /= 1000

	return temp, nil
}
