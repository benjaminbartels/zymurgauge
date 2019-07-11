package internal

import "time"

// Fermentation is a single instance of a fermentation of a beer
type Fermentation struct {
	ID            string     `json:"id"`
	Chamber       Chamber    `json:"chamber"`
	Beer          Beer       `json:"beer"`
	CurrentStep   int        `json:"currentStep,omitempty"`
	StartTime     *time.Time `json:"startTime,omitempty"`
	CompletedTime *time.Time `json:"completedTime,omitempty"`
	ModTime       time.Time  `json:"modTime"`
}
