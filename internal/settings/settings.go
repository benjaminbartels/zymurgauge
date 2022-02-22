package settings

import "time"

type Settings struct {
	BrewfatherAPIUserID string    `json:"brewfatherApiUserId"`
	BrewfatherAPIKey    string    `json:"brewfatherApiKey"`
	BrewfatherLogURL    string    `json:"brewfatherLogUrl"`
	TemperatureUnits    string    `json:"temperatureUnits"`
	ModTime             time.Time `json:"modTime"`
}
