package batch

import "github.com/benjaminbartels/zymurgauge/internal/brewfather"

type Summary struct {
	ID         string `json:"id"`
	Number     int    `json:"number"`
	RecipeName string `json:"recipeName"`
}

type Detail struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	Recipe Recipe `json:"recipe"`
}

type Recipe struct {
	Name            string       `json:"name"`
	Fermentation    Fermentation `json:"fermentation"`
	OriginalGravity float64      `json:"originalGravity"`
	FinalGravity    float64      `json:"finalGravity"`
}

type Fermentation struct {
	Name  string             `json:"name"`
	Steps []FermentationStep `json:"steps"`
}

type FermentationStep struct {
	Name        string  `json:"name"`
	Temperature float64 `json:"temperature"`
	Time        int     `json:"time"`
}

func ConvertSummaries(batches []brewfather.BatchSummary) []Summary {
	s := []Summary{}
	for i := 0; i < len(batches); i++ {
		s = append(s, convertSummary(batches[i]))
	}

	return s
}

func convertSummary(b brewfather.BatchSummary) Summary {
	return Summary{
		ID:         b.ID,
		Number:     b.BatchNo,
		RecipeName: b.Recipe.Name,
	}
}

func ConvertDetail(b *brewfather.BatchDetail) *Detail {
	return &Detail{
		ID:     b.ID,
		Number: b.BatchNo,
		Recipe: convertRecipe(b.Recipe),
	}
}

func convertRecipe(recipe brewfather.Recipe) Recipe {
	return Recipe{
		Name:            recipe.Name,
		Fermentation:    convertFermentation(recipe.Fermentation),
		OriginalGravity: recipe.Og,
		FinalGravity:    recipe.Fg,
	}
}

func convertFermentation(fermentation brewfather.Fermentation) Fermentation {
	return Fermentation{
		Name:  fermentation.Name,
		Steps: convertFermentationSteps(fermentation.Steps),
	}
}

func convertFermentationSteps(steps []brewfather.FermentationSteps) []FermentationStep {
	s := []FermentationStep{}
	for i := 0; i < len(steps); i++ {
		s = append(s, convertFermentationStep(steps[i]))
	}

	return s
}

func convertFermentationStep(step brewfather.FermentationSteps) FermentationStep {
	return FermentationStep{
		Name:        step.Type,
		Temperature: step.StepTemp,
		Time:        step.StepTime,
	}
}
