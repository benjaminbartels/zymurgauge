package gpio

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Thermostat struct {
	ThermometerID string        `json:"thermometerId"`
	CoolerGPIO    int           `json:"coolerGpio"`
	HeaterGPIO    int           `json:"heaterGpio"`
	Interval      time.Duration `json:"interval"`
	logger        logrus.Logger
	path          string
	target        *float64
	quit          chan bool
	isPolling     bool
}
