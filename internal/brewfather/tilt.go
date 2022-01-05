//nolint:tagliatelle
package brewfather

import "context"

type Service interface { // TODO: needs better name
	GetAll(ctx context.Context) ([]Batch, error)
	Get(ctx context.Context, id string) (*Batch, error)
	Log(ctx context.Context, log LogEntry) error
}

type Batch struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Fermentation Fermentation `json:"fermentation"`
}

type Fermentation struct {
	Name  string             `json:"name"`
	Steps []FermentationStep `json:"steps"`
}

type FermentationStep struct {
	Type       string  `json:"type"`
	ActualTime int64   `json:"actualTime"`
	StepTemp   float64 `json:"stepTemp"`
	StepTime   int     `json:"stepTime"`
}

type LogEntry struct {
	DeviceName          string `json:"name"`                    // Required field, first 15 characters used inID in Brewfather,
	BeerTemperature     string `json:"temp,omitempty"`          // 20.32,
	AuxilaryTemperature string `json:"aux_temp,omitempty"`      // 15.61, // Fridge Temp
	ExternalTemperature string `json:"ext_temp,omitempty"`      // 6.51, // Room Temp
	TemperatureUnit     string `json:"temp_unit,omitempty"`     // "C", // C, F, K
	Gravity             string `json:"gravity,omitempty"`       // 1.042,
	GravityUnit         string `json:"gravity_unit,omitempty"`  // "G", // G, P
	Preasrue            string `json:"pressure,omitempty"`      // 10,
	PreasrueUnit        string `json:"pressure_unit,omitempty"` // "PSI", // PSI, BAR, KPA
	Ph                  string `json:"ph,omitempty"`            // 4.12,
	BPM                 string `json:"bpm"`                     // 123, // Bubbles Per Minute
	Comment             string `json:"comment,omitempty"`       // "Hello World",
	Beer                string `json:"beer,omitempty"`          // "Pale Ale",
	Barttery            string `json:"battery,omitempty"`       // 4.98

}
