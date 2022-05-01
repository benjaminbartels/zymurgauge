package bluetooth_test

import (
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/platform/bluetooth"
	"github.com/stretchr/testify/assert"
)

func TestNewIBeacon(t *testing.T) {
	t.Parallel()

	bytes := []byte{
		76, 0, 2, 21, 164, 149, 187, 80, 197, 177, 75, 68, 181, 18, 19, 112, 240, 45, 116, 222, 0, 68, 3, 231, 197,
	}

	ibeacon, err := bluetooth.NewIBeacon(bytes)

	assert.Equal(t, "a495bb50c5b14b44b5121370f02d74de", ibeacon.GetUUID())
	assert.Equal(t, uint16(68), ibeacon.GetMajor())
	assert.Equal(t, uint16(999), ibeacon.GetMinor())
	assert.NoError(t, err)
}

func TestNewIBeaconError(t *testing.T) {
	t.Parallel()

	bytes := []byte{76, 0, 2, 21, 164}

	_, err := bluetooth.NewIBeacon(bytes)

	assert.ErrorIs(t, err, bluetooth.ErrInvalidManufacturerDataLength)
}
