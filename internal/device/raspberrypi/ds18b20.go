package raspberrypi

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/pkg/errors"
)

const (
	slave = "w1_slave"
	// defaultDevicePath is the location of the thermometer data on the file system.
	defaultDevicePath = "/sys/bus/w1/devices/"
	devicePrefix      = "28-"
)

// NewDs18b20 Create a new ds18b20 Thermometer.
func NewDs18b20(id string) (*Ds18b20, error) {
	return NewDs18b20WithDevicePath(id, defaultDevicePath)
}

// NewDs18b20WithDevicePath creates a new Thermometer using the given devicePath
// Usually used for testing.
func NewDs18b20WithDevicePath(id, devicePath string) (*Ds18b20, error) {
	_, err := os.Stat(path.Join(devicePath, id))
	if err != nil {
		return nil, errors.Wrap(err, "could not file information")
	}

	return &Ds18b20{
		ID:   id,
		path: devicePath,
	}, nil
}

// Ds18b20 is a GPIO 1-wire temperature probe.
type Ds18b20 struct {
	ID   string
	path string
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

var _ internal.ThermometerRepo = (*Ds18b20Repo)(nil)

type Ds18b20Repo struct {
	devicePath string
}

func NewDs18b20Repo() *Ds18b20Repo {
	return &Ds18b20Repo{devicePath: defaultDevicePath}
}

func (r *Ds18b20Repo) GetThermometerIDs() ([]string, error) {
	ids := []string{}

	dir, err := os.Open(defaultDevicePath)
	if err != nil {
		return ids, errors.Wrapf(err, "could not open %s", defaultDevicePath)
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return ids, errors.Wrapf(err, "could not open read directory names")
	}

	for _, name := range names {
		if strings.HasPrefix(name, devicePrefix) {
			ids = append(ids, name)
		}
	}

	return ids, nil
}
