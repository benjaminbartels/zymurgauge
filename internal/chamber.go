package internal

import "time"

type Chamber struct {
	MacAddress          string                 `json:"macAddress"`
	Name                string                 `json:"name"`
	Controller          *TemperatureController `json:"controller"`
	CurrentFermentation *Fermentation          `json:"currentFermentation"`
	ModTime             time.Time              `json:"modTime"`
}
