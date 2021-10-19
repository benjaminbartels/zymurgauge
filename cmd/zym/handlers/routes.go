package handlers

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/middleware"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/sirupsen/logrus"
)

const (
	chambersPath = "/chambers"
	recipesPath  = "/recipes"
	version      = "v1"
)

// NewAPI return a web.App with configured routes and handlers.
func NewAPI(chamberRepo chamber.Repo, recipeRepo brewfather.RecipeRepo, shutdown chan os.Signal,
	logger *logrus.Logger) http.Handler {
	chambersHandler := &ChambersHandler{
		Repo: chamberRepo,
	}

	recipesHandler := &RecipesHandler{
		Repo: recipeRepo,
	}

	app := web.NewApp(shutdown, middleware.RequestLogger(logger), middleware.Errors(logger))

	app.Register(http.MethodGet, version, chambersPath, chambersHandler.GetAll)
	app.Register(http.MethodGet, version, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Get)
	app.Register(http.MethodPost, version, chambersPath, chambersHandler.Save)
	app.Register(http.MethodDelete, version, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Delete)

	app.Register(http.MethodGet, version, recipesPath, recipesHandler.GetAll)
	app.Register(http.MethodGet, version, fmt.Sprintf("%s/:id", recipesPath), recipesHandler.Get)

	return app
}

func DebugMux() *http.ServeMux {
	debugMux := http.NewServeMux()
	debugMux.HandleFunc("/debug/pprof/", pprof.Index)
	debugMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	debugMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	debugMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	debugMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	debugMux.Handle("/debug/vars", expvar.Handler())

	return debugMux
}
