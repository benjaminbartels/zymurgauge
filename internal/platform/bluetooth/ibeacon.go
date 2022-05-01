package bluetooth

import (
	"encoding/binary"
	"encoding/hex"
)

const manufacurerDataMinLength = 25

type IBeacon struct {
	uuid  string
	major uint16
	minor uint16
}

func NewIBeacon(manufacurerData []byte) (*IBeacon, error) {
	if len(manufacurerData) < manufacurerDataMinLength {
		return nil, ErrInvalidManufacturerDataLength
	}

	return &IBeacon{
		uuid:  hex.EncodeToString(manufacurerData[4:20]),
		major: binary.BigEndian.Uint16(manufacurerData[20:22]),
		minor: binary.BigEndian.Uint16(manufacurerData[22:24]),
	}, nil
}

func (b *IBeacon) GetUUID() string {
	return b.uuid
}

func (b *IBeacon) GetMajor() uint16 {
	return b.major
}

func (b *IBeacon) GetMinor() uint16 {
	return b.minor
}
