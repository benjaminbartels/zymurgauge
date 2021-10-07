package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/benjaminbartels/zymurgauge/cmd/fermmon/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/thermometer"
	"github.com/kelseyhightower/envconfig"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"periph.io/x/periph/conn/onewire/onewirereg"
)

const (
	hours              = 24
	thermometerAddress = 2233785441622197288
	resolutionBits     = 10
)

type config struct {
	APIUserID     string  `required:"true"`
	APIKey        string  `required:"true"`
	ThermometerID uint64  `required:"true"`
	ChillerPIN    string  `required:"true"`
	HeaterPIN     string  `required:"true"`
	ChillerKp     float64 `required:"true"`
	ChillerKi     float64 `required:"true"`
	ChillerKd     float64 `required:"true"`
	HeaterKp      float64 `required:"true"`
	HeaterKi      float64 `required:"true"`
	HeaterKd      float64 `required:"true"`
	Debug         bool    `default:"false"`
}

func main() {
	var cfg config

	if err := envconfig.Process("fermmon", &cfg); err != nil {
		fmt.Println("Could not process env vars:", err)
		runtime.Goexit()
	}

	logger := logrus.New()

	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	bus, err := onewirereg.Open("")
	if err != nil {
		fmt.Println("Could not open 1-wire bus:", err)
		runtime.Goexit()
	}

	defer bus.Close()

	ds18b20, err := thermometer.NewDs18b20(bus, thermometerAddress, resolutionBits)
	if err != nil {
		fmt.Println("Could not create ds18b20 thermometer:", err)
		runtime.Goexit()
	}

	createFunc := CreateThermostat

	thermostat, err := createFunc(ds18b20, cfg.ChillerPIN, cfg.HeaterPIN, cfg.ChillerKp, cfg.ChillerKi,
		cfg.ChillerKd, cfg.HeaterKp, cfg.HeaterKi, cfg.HeaterKd, logger)
	if err != nil {
		fmt.Println("Could not create thermostat:", err)
		runtime.Goexit()
	}

	service := brewfather.New(brewfather.APIURL, cfg.APIUserID, cfg.APIKey)

	recipes, err := service.GetRecipes(context.Background())
	if err != nil {
		fmt.Println("Could not get Recipes:", err)
		runtime.Goexit()
	}

	id, err := runRecipesPrompt(recipes)
	if err != nil {
		fmt.Println("Could not run prompt for Recipes:", err)
		runtime.Goexit()
	}

	recipe, err := service.GetRecipe(context.Background(), id)
	if err != nil {
		fmt.Println("Could not run get Recipes:", err)
		runtime.Goexit()
	}

	startingStep, err := runFermentationStepsPrompt(recipe.Fermentation.Steps)
	if err != nil {
		fmt.Println("Could not run prompt for Fermentation Steps:", err)
		runtime.Goexit()
	}

	for i := startingStep; i < len(recipe.Fermentation.Steps); i++ {
		step := recipe.Fermentation.Steps[i]
		if err := thermostat.On(step.StepTemp); err != nil {
			fmt.Println("Could not turn thermostat on:", err)
			runtime.Goexit()
		}

		waitTimer := time.NewTimer(time.Duration(step.StepTime*hours) * time.Hour)
		<-waitTimer.C
		thermostat.Off()
	}
}

func runRecipesPrompt(recipes []brewfather.Recipe) (string, error) {
	prompt := promptui.Select{
		Label: "Select Recipe",
		Items: recipes,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U0001F37A {{ .Name | cyan }} ({{ .Style.Name | yellow }})",
			Inactive: "  {{ .Name | cyan }} ({{ .Style.Name | yellow }})",
			Selected: "\U0001F37A {{ .Name | cyan }}",
		},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return recipes[i].ID, nil
}

func runFermentationStepsPrompt(steps []brewfather.FermentationStep) (int, error) {
	prompt := promptui.Select{
		Label: "Select Fermentation Step",
		Items: steps,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U0001F37A {{ .Type | cyan }} ({{ .DisplayStepTemp | yellow }}, {{ .StepTime | yellow }} Days)",
			Inactive: "  {{ .Type | cyan }} ({{ .DisplayStepTemp | yellow }}°F, {{ .StepTime | yellow }} Days)",
			Selected: "\U0001F37A {{ .Type | cyan }} ({{ .DisplayStepTemp | yellow }}°F, {{ .StepTime | yellow }} Days)",
		},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	return i, nil
}
