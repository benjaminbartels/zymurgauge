package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/middleware"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
	"github.com/sirupsen/logrus"
)

const (
	chambersPath = "/chambers"
	recipesPath  = "/recipes"
)

// NewAPI return a web.App with configured routes and handlers.
func NewAPI(chamberRepo *storage.ChamberRepo, brewfatherClient *brewfather.Client, shutdown chan os.Signal,
	logger *logrus.Logger) http.Handler {
	chambersHandler := &Chambers{
		repo: chamberRepo,
	}

	recipesHandler := &Recipes{
		brewfatherClient: brewfatherClient,
	}

	// TODO: Allow for Versions

	app := web.NewApp(shutdown, middleware.RequestLogger(logger), middleware.Errors(logger))

	app.Register(http.MethodGet, chambersPath, chambersHandler.GetAll)
	app.Register(http.MethodGet, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Get)
	app.Register(http.MethodPost, chambersPath, chambersHandler.Save)
	app.Register(http.MethodDelete, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Delete)
	app.Register(http.MethodGet, recipesPath, recipesHandler.GetAll)
	app.Register(http.MethodGet, fmt.Sprintf("%s/:id", recipesPath), recipesHandler.Get)

	return app
}
