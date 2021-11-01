package batch

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
