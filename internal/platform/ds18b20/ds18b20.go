package ds18b20

import (
	"bufio"
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/benjaminbartels/zymurgauge/internal/platform/safeclose"
)

const (
	slave = "w1_slave"
)

// defaultDevicePath is the location of the thermometer data on the file system
const defaultDevicePath = "/sys/bus/w1/devices/"

// New Create a new ds18b20 Thermometer
func New(id string) (*Thermometer, error) {
	return NewWithDevicePath(id, defaultDevicePath)
}

// NewWithDevicePath creates a new Thermometer using the given devicePath
// Usually used for testing
func NewWithDevicePath(id, devicePath string) (*Thermometer, error) {
	_, err := os.Stat(path.Join(devicePath, id))
	if err != nil {
		return nil, err
	}

	return &Thermometer{
		ID:   id,
		path: devicePath,
	}, nil
}

// Thermometer is a GPIO temperature probe
type Thermometer struct {
	ID   string
	path string
}

// ReadTemperature read the current temperature of the Thermometer
func (t *Thermometer) Read() (*float64, error) {
	file, err := os.Open(path.Join(t.path, t.ID, slave))
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
