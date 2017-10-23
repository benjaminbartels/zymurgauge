package internal

import "time"

// TemperatureChange is the event that represents a change in temperature
type TemperatureChange struct {
	ID             uint64    `json:"id"`
	FermentationID uint64    `json:"fermentationId"`
	InsertTime     time.Time `json:"insertTime"`
	Chamber        string    `json:"chamber"`
	Beer           string    `json:"beer"`
	Thermometer    string    `json:"thermometer"`
	Temperature    float64   `json:"temperature"`
}

// type ThermostatEvent struct {
// 	FermentationEvent
// 	Chamber string          `json:"chamber"`
// 	Beer    string          `json:"beer"`
// 	State   ThermostatState `json:"state"`
// }
