package internal

import "time"

// Fermentation is a single instance of a fermentation of a beer
type Fermentation struct {
	ID            uint64    `json:"id"`
	Beer          Beer      `json:"beer"`
	CurrentStep   int       `json:"currentStep,omitempty"`
	StartTime     time.Time `json:"startTime,omitempty"`
	CompletedTime time.Time `json:"completedTime,omitempty"`
	ModTime       time.Time `json:"modTime"`
}
