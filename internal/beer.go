package internal

import (
	"time"
)

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
