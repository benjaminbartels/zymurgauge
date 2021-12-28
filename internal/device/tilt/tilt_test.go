package tilt_test

import (
	"testing"

	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
	"github.com/stretchr/testify/assert"
)

func TestGetTemperatureError(t *testing.T) {
	t.Parallel()

	tilt := &tilt.Tilt{}
	temp, err := tilt.GetTemperature()

	assert.Equal(t, 0.0, temp)

	// TODO: Waiting on PR for ErrorContains(): https://github.com/stretchr/testify/pull/1022
	assert.Contains(t, err.Error(), "underlying IBeacon is nil")
}

func TestGetSpecificGravityError(t *testing.T) {
	t.Parallel()

	tilt := &tilt.Tilt{}
	temp, err := tilt.GetSpecificGravity()

	assert.Equal(t, 0.0, temp)

	// TODO: Waiting on PR for ErrorContains(): https://github.com/stretchr/testify/pull/1022
	assert.Contains(t, err.Error(), "underlying IBeacon is nil")
}
