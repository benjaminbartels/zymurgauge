package internal

import "time"

// Chamber represents an insulated box (fridge) with internal heating/cooling elements that reacts to changes in
// monitored temperatures, by correcting small deviations from your desired fermentation temperature.
type Chamber struct {
	MacAddress            string      `json:"macAddress"`
	Name                  string      `json:"name"`
	Thermostat            *Thermostat `json:"thermostat"`
	CurrentFermentationID *uint64     `json:"currentFermentationId,omitempty"`
	ModTime               time.Time   `json:"modTime"`
}
