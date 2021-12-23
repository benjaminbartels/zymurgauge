package onewire

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	slave = "w1_slave"
	// defaultDevicePath is the location of the thermometer data on the file system.
	DefaultDevicePath = "/sys/bus/w1/devices/"
	devicePrefix      = "28-*"
)

type OptionsFunc func(*Ds18b20)

// Ds18b20 is a GPIO 1-wire temperature probe.
type Ds18b20 struct {
	ID   string
	path string
}

func SetDevicePath(path string) OptionsFunc {
	return func(d *Ds18b20) {
		d.path = path
	}
}

// NewDs18b20 Create a new ds18b20 Thermometer.
func NewDs18b20(id string, options ...OptionsFunc) (*Ds18b20, error) {
	d := &Ds18b20{
		ID:   id,
		path: DefaultDevicePath,
	}

	for _, option := range options {
		option(d)
	}

	_, err := os.Stat(path.Join(d.path, id))
	if err != nil {
		return nil, errors.Wrapf(err, "could not get file information for %s", d.path)
	}

	return d, nil
}

// ReadTemperature read the current temperature of the Thermometer.
func (d *Ds18b20) GetTemperature() (float64, error) {
	file, err := os.Open(path.Join(d.path, d.ID, slave))
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

func GetThermometerIDs(devicePath string) ([]string, error) {
	p := path.Join(devicePath, devicePrefix)

	filenames, err := filepath.Glob(p)
	if err != nil {
		return nil, errors.Wrapf(err, "could matching files matching %s", p)
	}

	return filenames, nil
}
