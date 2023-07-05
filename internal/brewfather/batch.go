//nolint:tagliatelle // this file was manual generated and includes minor changes
package brewfather

type BatchSummary struct {
	ID      string `json:"_id"`
	BatchNo int    `json:"batchNo"`
	Recipe  struct {
		Name string `json:"name"`
	} `json:"recipe"`
}

type BatchDetail struct {
	ID      string `json:"_id"`
	BatchNo int    `json:"batchNo"`
	Recipe  Recipe `json:"recipe"`
}

type Recipe struct {
	Name         string       `json:"name"`
	Fermentation Fermentation `json:"fermentation"`
	Fg           float64      `json:"fg"`
	Og           float64      `json:"og"`
}

type Fermentation struct {
	ID    string              `json:"_id"`
	Name  string              `json:"name"`
	Steps []FermentationSteps `json:"steps"`
}

type FermentationSteps struct {
	Type     string  `json:"type"`
	StepTemp float64 `json:"stepTemp"`
	StepTime int     `json:"stepTime"`
}
