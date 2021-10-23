package raspberrypi

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/pkg/errors"
	"periph.io/x/periph/conn/onewire"
	"periph.io/x/periph/devices/ds18b20"
	"periph.io/x/periph/experimental/host/netlink"
)

var _ device.Thermometer = (*Ds18b20)(nil)

type Ds18b20 struct {
	dev *ds18b20.Dev
}

func NewDs18b20(bus onewire.Bus, id string, resolutionBits int) (*Ds18b20, error) {
	address, err := getAddress(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not get ds18b20 address")
	}

	dev, err := ds18b20.New(bus, address, resolutionBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not create ds18b20 thermometer")
	}

	return &Ds18b20{
		dev: dev,
	}, nil
}

func (d *Ds18b20) GetTemperature() (float64, error) {
	temp, err := d.dev.LastTemp()
	if err != nil {
		return 0, errors.Wrap(err, "could not read ds18b20 thermometer")
	}

	return temp.Celsius(), nil
}

func GetThermometerIDs(oneBus *netlink.OneWire) ([]string, error) {
	ids := []string{}

	addresses, err := oneBus.Search(false)
	if err != nil {
		return ids, errors.Wrap(err, "could not search bus")
	}

	for _, address := range addresses {
		ids = append(ids, getID(address))
	}

	return ids, nil
}

func getID(address onewire.Address) string {
	bytes := make([]byte, 8) // nolint: gomnd // remove when https://github.com/tommy-muehle/go-mnd/pull/29 is merged
	binary.LittleEndian.PutUint64(bytes, uint64(address))

	return fmt.Sprintf("%s-%s", hex.EncodeToString(bytes[0:1]), hex.EncodeToString(reverse(bytes[1:7])))
}

func getAddress(id string) (onewire.Address, error) {
	s := strings.Split(id, "-")

	family, err := hex.DecodeString(s[0])
	if err != nil {
		return 0, errors.Wrap(err, "could not decode ds18b20 family code")
	}

	uid, err := hex.DecodeString(s[1])
	if err != nil {
		return 0, errors.Wrap(err, "could not decode unique device id")
	}

	bytes := make([]byte, 8) // nolint: gomnd // remove when https://github.com/tommy-muehle/go-mnd/pull/29 is merged
	copy(bytes, append(family, reverse(uid)...))

	return onewire.Address(binary.LittleEndian.Uint64(bytes)), nil
}

func reverse(bytes []byte) []byte {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	return bytes
}
