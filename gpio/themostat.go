package gpio

import (
	"time"

	"github.com/sirupsen/logrus"
)

type state int

const (
	OFF state = 1 + iota
	COOLING
	HEATING
)

type Thermostat struct {
	ThermometerID string        `json:"thermometerId"`
	CoolerGPIO    *int          `json:"coolerGpio"`
	HeaterGPIO    *int          `json:"heaterGpio"`
	Interval      time.Duration `json:"interval"`
	Logger        *logrus.Logger
	path          string
	target        *float64
	quit          chan bool
	isPolling     bool
	state         state
}
