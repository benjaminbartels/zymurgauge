package tilt

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/pkg/errors"
)

const minLength = 25

var ErrInvalidManufacturerDataLength = errors.New("manufacurerData length is less that 25")

type IBeacon struct {
	UUID  string
	Major uint16
	Minor uint16
}

func NewIBeacon(manufacurerData []byte) (*IBeacon, error) {
	if len(manufacurerData) < minLength {
		return nil, ErrInvalidManufacturerDataLength
	}

	return &IBeacon{
		UUID:  hex.EncodeToString(manufacurerData[4:20]), // TODO: remove preamble
		Major: binary.BigEndian.Uint16(manufacurerData[20:22]),
		Minor: binary.BigEndian.Uint16(manufacurerData[22:24]),
	}, nil
}
