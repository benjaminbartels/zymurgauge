package storage

import "time"

type ChamberRepo interface {
	GetAll() ([]Chamber, error)
	Get(id string) (*Chamber, error)
	Save(c *Chamber) error
	Delete(id string) error
} // TODO: Move this?

// Chamber represents an insulated box (fridge) with internal heating/cooling elements that reacts to changes in
// monitored temperatures, by correcting small deviations from your desired fermentation temperature.
type Chamber struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	ThermometerID string    `json:"thermometerId"`
	ChillerPIN    string    `json:"chillerPin"`
	HeaterPIN     string    `json:"heaterPin"`
	ChillerKp     float64   `json:"chillerKp"`
	ChillerKi     float64   `json:"chillerKi"`
	ChillerKd     float64   `json:"chillerKd"`
	HeaterKp      float64   `json:"heaterKp"`
	HeaterKi      float64   `json:"heaterKi"`
	HeaterKd      float64   `json:"heaterKd"`
	ModTime       time.Time `json:"modTime"`
}

// TODO: Move this to common location?
