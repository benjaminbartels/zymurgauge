package main

import (
	"github.com/benjaminbartels/zymurgauge/internal/device/raspberrypi"
)

func createThermometerRepo() *raspberrypi.Ds18b20Repo {
	return raspberrypi.NewDs18b20Repo()
}
