package ds18b20

import (
	"bufio"
	"errors"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/benjaminbartels/zymurgauge/internal/platform/safeclose"
	"github.com/spf13/afero"
)

const (
	devicePath = "/sys/bus/w1/devices/"
	prefix     = "28-"
	slave      = "w1_slave"
)

// ThermometerGroup is the group of Thermometers on the system
type ThermometerGroup struct {
	fs afero.Fs
}

// New returns a ThermometerGroup using the specified FileSystem
func New(fs afero.Fs) *ThermometerGroup {
	return &ThermometerGroup{
		fs: fs,
	}
}

// GetThermometers returns a list of all Thermometers on the bus
func (t *ThermometerGroup) GetThermometers() ([]Thermometer, error) {
	var thermometers []Thermometer

	dir, err := t.fs.Open(devicePath)
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
				ID:    info.Name(),
				fs:    t.fs,
				mutex: &sync.Mutex{},
			}

			thermometers = append(thermometers, term)
		}
	}
	return thermometers, nil
}

// GetThermometer returns a Thermometer by id
func (t *ThermometerGroup) GetThermometer(id string) (*Thermometer, error) {

	_, err := t.fs.Stat(path.Join(devicePath, id))
	if err != nil {
		return nil, err
	}

	return &Thermometer{
		ID:    id,
		fs:    t.fs,
		mutex: &sync.Mutex{},
	}, nil
}

// Thermometer is a GPIO temperature probe
type Thermometer struct {
	ID    string
	fs    afero.Fs
	mutex *sync.Mutex
}

// ReadTemperature read the current temperature of the Thermometer
func (t *Thermometer) ReadTemperature() (*float64, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	file, err := t.fs.Open(path.Join(devicePath, t.ID, slave))
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
