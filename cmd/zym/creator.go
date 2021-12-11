//go:build !linux || !arm
// +build !linux !arm

package main

import "github.com/benjaminbartels/zymurgauge/internal/test/mocks"

func createThermometerRepo() *mocks.ThermometerRepo {
	ids := []string{"fake_thermometer_1", "fake_thermometer_2", "fake_thermometer_3"}
	repoMock := &mocks.ThermometerRepo{}
	repoMock.On("GetThermometerIDs").Return(ids, nil)

	return repoMock
}
