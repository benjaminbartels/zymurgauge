package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/middleware"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/settings"
	"github.com/benjaminbartels/zymurgauge/ui"
	"github.com/sirupsen/logrus"
)

const (
	chambersPath     = "/chambers"
	thermometersPath = "/thermometers"
	batchesPath      = "/batches"
	settingsPath     = "/settings"
	version          = "v1"
)

func NewApp(chamberManager chamber.Controller, devicePath string,
	service brewfather.Service, settingsRepo settings.Repo, updateChan chan settings.Settings,
	uiFileReader web.FileReader, shutdown chan os.Signal, logger *logrus.Logger) *web.App {
	api := web.NewAPI(shutdown,
		middleware.RequestLogger(logger),
		middleware.Errors(logger),
		middleware.Cors())

	chambersHandler := &ChambersHandler{
		ChamberController: chamberManager,
		Logger:            logger,
	}

	api.Register(http.MethodGet, version, chambersPath, chambersHandler.GetAll)
	api.Register(http.MethodGet, version, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Get)
	api.Register(http.MethodPost, version, chambersPath, chambersHandler.Save)
	api.Register(http.MethodDelete, version, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Delete)
	api.Register(http.MethodPost, version, fmt.Sprintf("%s/:id/start", chambersPath), chambersHandler.Start)
	api.Register(http.MethodPost, version, fmt.Sprintf("%s/:id/stop", chambersPath), chambersHandler.Stop)

	batchesHandler := &BatchesHandler{
		Service: service,
	}

	api.Register(http.MethodGet, version, batchesPath, batchesHandler.GetAll)
	api.Register(http.MethodGet, version, fmt.Sprintf("%s/:id", batchesPath), batchesHandler.Get)

	thermometersHandler := &ThermometersHandler{
		DevicePath: devicePath,
	}

	api.Register(http.MethodGet, version, thermometersPath, thermometersHandler.GetAll)

	settingsHandler := &SettingsHandler{
		SettingsRepo: settingsRepo,
		UpdateChan:   updateChan,
	}

	api.Register(http.MethodGet, version, settingsPath, settingsHandler.Get)
	api.Register(http.MethodPost, version, settingsPath, settingsHandler.Save)

	app := web.NewApp(api, ui.FS, logger)

	return app
}
