package settings

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/auth"
)

type Settings struct {
	AppSettings
	auth.Credentials
	ModTime time.Time `json:"modTime"`
}

type AppSettings struct {
	TemperatureUnits    string `json:"temperatureUnits"`
	AuthSecret          string `json:"authSecret"`
	BrewfatherAPIUserID string `json:"brewfatherApiUserId,omitempty"`
	BrewfatherAPIKey    string `json:"brewfatherApiKey,omitempty"`
	BrewfatherLogURL    string `json:"brewfatherLogUrl,omitempty"`
	InfluxDBURL         string `json:"influxDbUrl,omitempty"`
	InfluxDBReadToken   string `json:"influxDbReadToken,omitempty"`
	StatsDAddress       string `json:"statsDAddress,omitempty"`
}
