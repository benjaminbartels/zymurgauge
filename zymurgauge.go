package zymurgauge

import "time"

// Client creates a connection to the services
type Client interface {
	ChamberService() ChamberService
	FermentationService() FermentationService
	BeerService() BeerService
}

// ChamberService is used to manage fermentation controller
type ChamberService interface {
	Get(mac string) (*Chamber, error)
	Save(f *Chamber) error
	Subscribe(mac string, ch chan Chamber) error
	Unsubscribe(mac string)
}

// FermentationService is used to manage fermentations
type FermentationService interface {
	Get(id uint64) (*Fermentation, error)
	Save(f *Fermentation) error
	// ToDo: Move to its own service
	//LogEvent(fermentationID uint64, event FermentationEvent) error
}

// BeerService is used to manage beers
type BeerService interface {
	Get(id uint64) (*Beer, error)
	Save(b *Beer) error
}

type TemperatureController interface {
	SetTemperature(t *float64) error
}

type Chamber struct {
	MacAddress          string                `json:"macAddress"`
	Name                string                `json:"name"`
	Controller          TemperatureController `json:"controller"`
	CurrentFermentation *Fermentation         `json:"currentFermentation"`
	ModTime             time.Time             `json:"modTime"`
}

// Fermentation is a single instance of a fermentation of a beer
type Fermentation struct {
	ID            uint64    `json:"id"`
	Beer          Beer      `json:"beer"`
	CurrentStep   int       `json:"currentStep"`
	StartTime     time.Time `json:"startTime"`
	CompletedTime time.Time `json:"completedTime"`
	ModTime       time.Time `json:"modTime"`
}

// Beer represents details of a beer to be fermented
type Beer struct {
	ID       uint64             `json:"id"`
	Name     string             `json:"name"`
	Style    string             `json:"style"`
	Schedule []FermentationStep `json:"schedule"`
	ModTime  time.Time          `json:"modTime"`
}

// FermentationStep represents the fermentation details for a single step of a fermentation schedule
type FermentationStep struct {
	Order      int           `json:"order"`
	TargetTemp float64       `json:"targetTemp"`
	Duration   time.Duration `json:"duration"`
}

// FermentationEvent is an event that occurred during the fermentation of a beer.  Types include
// Start, Stop, TemperatureReading, StepChange, Complete, CoolingOn, CoolingOff, HeatingOn, HeatingOff
type FermentationEvent struct {
	FermentationID uint64    `json:"fermentationId"`
	Time           time.Time `json:"eventTime"`
	Type           string    `json:"type"`
	Value          string    `json:"value"`
}
