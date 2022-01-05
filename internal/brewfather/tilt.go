//nolint:tagliatelle
package brewfather

import "context"

type Service interface { // TODO: needs better name
	GetAll(ctx context.Context) ([]Batch, error)
	Get(ctx context.Context, id string) (*Batch, error)
	LogTilt(ctx context.Context, log TiltLogEntry) error
}

type TiltLogEntry struct {
	// Timepoint       string `json:"Timepoint"` //TODO: is this needed
	Temperature     string `json:"Temp"`
	SpecificGravity string `json:"SG"`
	Beer            string `json:"Beer"`
	Color           string `json:"Color"`
	Comment         string `json:"Comment"`
}
