package settings

import "time"

type Settings struct {
	TemperatureUnits    string    `json:"temperatureUnits"`
	AuthSecret          string    `json:"authSecret"`
	BrewfatherAPIUserID string    `json:"brewfatherApiUserId,omitempty"`
	BrewfatherAPIKey    string    `json:"brewfatherApiKey,omitempty"`
	BrewfatherLogURL    string    `json:"brewfatherLogUrl,omitempty"`
	InfluxDBURL         string    `json:"influxDbUrl,omitempty"`
	StatsDAddress       string    `json:"statsDAddress,omitempty"`
	ModTime             time.Time `json:"modTime"`
}
