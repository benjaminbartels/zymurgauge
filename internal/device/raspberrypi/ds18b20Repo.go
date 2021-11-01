package raspberrypi

import (
	"os"
	"strings"

	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/pkg/errors"
)

var _ device.ThermometerRepo = (*Ds18b20Repo)(nil)

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
