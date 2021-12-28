package tilt_test

import (
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
	"github.com/stretchr/testify/assert"
)

func TestNewIBeacon(t *testing.T) {
	t.Parallel()

	bytes := []byte{
		76, 0, 2, 21, 164, 149, 187, 80, 197, 177, 75, 68, 181, 18, 19, 112, 240, 45, 116, 222, 0, 68, 3, 231, 197,
	}

	ibeacon, err := tilt.NewIBeacon(bytes)

	assert.Equal(t, "a495bb50c5b14b44b5121370f02d74de", ibeacon.UUID)
	assert.Equal(t, uint16(68), ibeacon.Major)
	assert.Equal(t, uint16(999), ibeacon.Minor)
	assert.NoError(t, err)
}

func TestNewIBeaconError(t *testing.T) {
	t.Parallel()

	bytes := []byte{76, 0, 2, 21, 164}

	_, err := tilt.NewIBeacon(bytes)

	assert.ErrorIs(t, err, tilt.ErrInvalidManufacturerDataLength)
}
