package gpio

import (
	"github.com/benjaminbartels/zymurgauge"
	"github.com/sirupsen/logrus"
)

type state int

const (
	OFF state = 1 + iota
	COOLING
	HEATING
)

type Thermostat struct {
	zymurgauge.TemperatureController
	Logger    *logrus.Logger
	path      string
	target    *float64
	quit      chan bool
	isPolling bool
	state     state
}
