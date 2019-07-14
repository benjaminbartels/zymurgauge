package controller_test

import (
	golog "log"
	"os"
	"testing"
	"time"

	"github.com/benjaminbartels/zymurgauge/cmd/fermmon/controller"
	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/client"
	"github.com/felixge/pidctrl"
)

const (
	mac = "00:11:22:33:44:55"
)

func TestController(t *testing.T) {
	t.Run("NoFermentationToNewFermentation", testNoFermentationToNewFermentation)
	t.Run("ReplaceFermentation", testReplaceFermentation)
	t.Run("FermentationToNoFermentation", testFermentationToNoFermentation)
}

func testNoFermentationToNewFermentation(t *testing.T) {

	fermentation := createFermentation("1")
	chamber := createChamber("")

	chamberResource := &chamberResourceMock{}
	chamberResource.GetFn = func(mac string) (*internal.Chamber, error) {
		return chamber, nil
	}

	fermentationResource := &fermentationResourceMock{}
	fermentationResource.GetFn = func(id string) (*internal.Fermentation, error) {
		return fermentation, nil
	}

	fermentationResource.SaveTemperatureChangeFn = func(*internal.TemperatureChange) error {
		return nil
	}

	ctl := createController(mac, chamberResource, fermentationResource)

	ctl.ConfigFunc = controller.ConfigureStubThermostat

	ctl.Poll()

	if !chamberResource.GetInvoked {
		t.Fatal("Chamber.Get not invoked")
	}

	if fermentationResource.GetInvoked {
		t.Fatal("Fermentation.Get invoked")
	}

	if ctl.Chamber.CurrentFermentationID != "" {
		t.Fatal("CurrentFermentationID is not empty")
	}

	if ctl.Fermentation != nil {
		t.Fatal("Fermentation is not nil")
	}

	chamber = createChamber("1")
	chamberResource.GetInvoked = false

	ctl.Poll()

	if !chamberResource.GetInvoked {
		t.Fatal("Chamber.Get not invoked")
	}

	if !fermentationResource.GetInvoked {
		t.Fatal("Fermentation.Get not invoked")
	}

	if ctl.Chamber.CurrentFermentationID != chamber.CurrentFermentationID {
		t.Fatal("CurrentFermentationID is not", chamber.CurrentFermentationID)
	}

	if ctl.Fermentation == nil {
		t.Fatal("Fermentation is nil")
	}
}

func testReplaceFermentation(t *testing.T) {

	fermentation1 := createFermentation("1")
	fermentation2 := createFermentation("2")
	chamber := createChamber("1")

	chamberResource := &chamberResourceMock{}
	chamberResource.GetFn = func(mac string) (*internal.Chamber, error) {
		return chamber, nil
	}

	fermentationResource := &fermentationResourceMock{}
	fermentationResource.GetFn = func(id string) (*internal.Fermentation, error) {

		if id == "1" {
			return fermentation1, nil
		} else if id == "2" {
			return fermentation2, nil
		}

		return nil, nil
	}

	fermentationResource.SaveTemperatureChangeFn = func(*internal.TemperatureChange) error {
		return nil
	}

	ctl := createController(mac, chamberResource, fermentationResource)

	ctl.ConfigFunc = controller.ConfigureStubThermostat

	ctl.Poll()

	if !chamberResource.GetInvoked {
		t.Fatal("Chamber.Get not invoked")
	}

	if !fermentationResource.GetInvoked {
		t.Fatal("Fermentation.Get not invoked")
	}

	if ctl.Chamber.CurrentFermentationID != chamber.CurrentFermentationID {
		t.Fatal("CurrentFermentationID is not", chamber.CurrentFermentationID)
	}

	if ctl.Fermentation == nil {
		t.Fatal("Fermentation is nil")
	}

	chamber = createChamber("2")
	chamberResource.GetInvoked = false
	fermentationResource.GetInvoked = false

	ctl.Poll()

	if !chamberResource.GetInvoked {
		t.Fatal("Chamber.Get not invoked")
	}

	if !fermentationResource.GetInvoked {
		t.Fatal("Fermentation.Get not invoked")
	}

	if ctl.Chamber.CurrentFermentationID != chamber.CurrentFermentationID {
		t.Fatal("CurrentFermentationID is not", chamber.CurrentFermentationID)
	}

	if ctl.Fermentation == nil {
		t.Fatal("Fermentation is nil")
	}
}

func testFermentationToNoFermentation(t *testing.T) {

	fermentation := createFermentation("1")
	chamber := createChamber("1")

	chamberResource := &chamberResourceMock{}
	chamberResource.GetFn = func(mac string) (*internal.Chamber, error) {
		return chamber, nil
	}

	fermentationResource := &fermentationResourceMock{}
	fermentationResource.GetFn = func(id string) (*internal.Fermentation, error) {
		return fermentation, nil
	}

	fermentationResource.SaveTemperatureChangeFn = func(*internal.TemperatureChange) error {
		return nil
	}

	ctl := createController(mac, chamberResource, fermentationResource)

	ctl.ConfigFunc = controller.ConfigureStubThermostat

	ctl.Poll()

	if !chamberResource.GetInvoked {
		t.Fatal("Chamber.Get not invoked")
	}

	if !fermentationResource.GetInvoked {
		t.Fatal("Fermentation.Get not invoked")
	}

	if ctl.Chamber.CurrentFermentationID != chamber.CurrentFermentationID {
		t.Fatal("CurrentFermentationID is not", chamber.CurrentFermentationID)
	}

	if ctl.Fermentation == nil {
		t.Fatal("Fermentation is nil")
	}

	chamber = createChamber("")
	chamberResource.GetInvoked = false
	fermentationResource.GetInvoked = false

	ctl.Poll()

	if !chamberResource.GetInvoked {
		t.Fatal("Chamber.Get not invoked")
	}

	if fermentationResource.GetInvoked {
		t.Fatal("Fermentation.Get invoked")
	}

	if ctl.Chamber.CurrentFermentationID != "" {
		t.Fatal("CurrentFermentationID is not empty")
	}

	if ctl.Fermentation != nil {
		t.Fatal("Fermentation is not nil")
	}
}

func createChamber(fermentationID string) *internal.Chamber {
	return &internal.Chamber{
		MacAddress: mac,
		Name:       "Chamber " + mac,
		Thermostat: &internal.Thermostat{
			ThermometerID: "1a2b3c4d",
			ChillerPin:    "1",
			HeaterPin:     "2",
		},
		CurrentFermentationID: fermentationID,
	}
}

func createFermentation(id string) *internal.Fermentation {
	var now = time.Now()

	return &internal.Fermentation{
		ID:          id,
		CurrentStep: 1,
		StartTime:   &now,
		ModTime:     now,
		Chamber: internal.Chamber{
			MacAddress: mac,
			Name:       "Chamber " + mac,
			Thermostat: &internal.Thermostat{
				ThermometerID: "1a2b3c4d",
				ChillerPin:    "1",
				HeaterPin:     "2",
			},
		},
		Beer: internal.Beer{
			ID:    "1",
			Name:  "My Beer",
			Style: "My Style",
			Schedule: []internal.FermentationStep{
				{
					Order:      1,
					TargetTemp: 16,
					Duration:   24 * time.Hour,
				},
			},
		},
	}
}

func createController(mac string, c client.ChamberProvider, f client.FermentationProvider) *controller.Controller {

	pid := pidctrl.NewPIDController(1, 1, 0) // ToDo: get from env vars

	logger := golog.New(os.Stderr, "", golog.LstdFlags)

	return controller.New(mac, pid, c, f, logger,
		internal.MinimumChill(3*time.Minute),
		internal.MinimumHeat(3*time.Minute),
		internal.Interval(10*time.Second),
		internal.Logger(logger))
}
