//nolint:tagliatelle // this file was manual generated and includes minor changes
package brewfather

import "context"

type Service interface {
	GetAllBatchSummaries(ctx context.Context) ([]BatchSummary, error)
	GetBatchDetail(ctx context.Context, id string) (*BatchDetail, error)
	Log(ctx context.Context, log LogEntry) error
}

type LogEntry struct {
	DeviceName           string `json:"name"`                    // Required, first 15 characters used in ID in Brewfather
	BeerTemperature      string `json:"temp,omitempty"`          // 20.32,
	AuxiliaryTemperature string `json:"aux_temp,omitempty"`      // 15.61, // Fridge Temp
	ExternalTemperature  string `json:"ext_temp,omitempty"`      // 6.51, // Room Temp
	TemperatureUnit      string `json:"temp_unit,omitempty"`     // "C", // C, F, K
	Gravity              string `json:"gravity,omitempty"`       // 1.042,
	GravityUnit          string `json:"gravity_unit,omitempty"`  // "G", // G, P
	Pressure             string `json:"pressure,omitempty"`      // 10,
	PressureUnit         string `json:"pressure_unit,omitempty"` // "PSI", // PSI, BAR, KPA
	Ph                   string `json:"ph,omitempty"`            // 4.12,
	BPM                  string `json:"bpm"`                     // 123, // Bubbles Per Minute
	Comment              string `json:"comment,omitempty"`       // "Hello World",
	Beer                 string `json:"beer,omitempty"`          // "Pale Ale",
	Battery              string `json:"battery,omitempty"`       // 4.98
}
