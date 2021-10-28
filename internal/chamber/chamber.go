package chamber

import (
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/device"
)

// Chamber represents an insulated box (fridge) with internal heating/cooling elements that reacts to changes in
// monitored temperatures, by correcting small deviations from your desired fermentation temperature.
type Chamber struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	ThermometerAddress uint64    `json:"thermometerAddress"`
	ChillerPIN         string    `json:"chillerPin"`
	HeaterPIN          string    `json:"heaterPin"`
	ChillerKp          float64   `json:"chillerKp"`
	ChillerKi          float64   `json:"chillerKi"`
	ChillerKd          float64   `json:"chillerKd"`
	HeaterKp           float64   `json:"heaterKp"`
	HeaterKi           float64   `json:"heaterKi"`
	HeaterKd           float64   `json:"heaterKd"`
	ModTime            time.Time `json:"modTime"`
	PIDController      device.TemperatureController
}
