package tilt

import (
	"encoding/binary"
	"encoding/hex"
)

type IBeacon struct {
	UUID  string
	Major uint16
	Minor uint16
}

func NewIBeacon(manufacurerData []byte) *IBeacon {
	return &IBeacon{
		UUID:  hex.EncodeToString(manufacurerData[4:20]), // TODO: remove preamble
		Major: binary.BigEndian.Uint16(manufacurerData[20:22]),
		Minor: binary.BigEndian.Uint16(manufacurerData[22:24]),
	}
}
