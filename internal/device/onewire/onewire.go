package onewire

import (
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

type Prefix string

const (
	Ds18b20Prefix Prefix = "28-"

	DefaultDevicePath = "/sys/bus/w1/devices/"
)

func GetIDs(devicePath string, prefix Prefix) ([]string, error) {
	p := path.Join(devicePath, string(prefix)+"*")

	filenames, err := filepath.Glob(p)
	if err != nil {
		return nil, errors.Wrapf(err, "could matching files matching %s", p)
	}

	result := []string{}

	for i := range filenames {
		result = append(result, filepath.Base(filenames[i]))
	}

	return result, nil
}
