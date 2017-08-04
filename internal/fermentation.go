package internal

import "time"

// Fermentation is a single instance of a fermentation of a beer
type Fermentation struct {
	ID            uint64    `json:"id"`
	Beer          Beer      `json:"beer"`
	CurrentStep   int       `json:"currentStep"`
	StartTime     time.Time `json:"startTime"`
	CompletedTime time.Time `json:"completedTime"`
	ModTime       time.Time `json:"modTime"`
}
