package handlers

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/middleware"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	uiweb "github.com/benjaminbartels/zymurgauge/web"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

const (
	chambersPath     = "/chambers"
	thermometersPath = "/thermometers"
	batchesPath      = "/batches"
	version          = "v1"
	uiDir            = "web/build"
	base             = "/ui"
)

// NewAPI return a web.App with configured routes and handlers.
func NewAPI(chamberManager chamber.Controller, devicePath string, service brewfather.Service, uiFiles uiweb.FileReader,
	shutdown chan os.Signal, logger *logrus.Logger) http.Handler {
	app := web.NewApp(shutdown, middleware.RequestLogger(logger), middleware.Errors(logger))

	chambersHandler := &ChambersHandler{
		ChamberController: chamberManager,
		Logger:            logger,
	}

	app.Register(http.MethodGet, version, chambersPath, chambersHandler.GetAll)
	app.Register(http.MethodGet, version, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Get)
	app.Register(http.MethodPost, version, chambersPath, chambersHandler.Save)
	app.Register(http.MethodDelete, version, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Delete)
	app.Register(http.MethodPost, version, fmt.Sprintf("%s/:id/start", chambersPath), chambersHandler.Start)
	app.Register(http.MethodPost, version, fmt.Sprintf("%s/:id/stop", chambersPath), chambersHandler.Stop)

	batchesHandler := &BatchesHandler{
		Service: service,
	}

	app.Register(http.MethodGet, version, batchesPath, batchesHandler.GetAll)
	app.Register(http.MethodGet, version, fmt.Sprintf("%s/:id", batchesPath), batchesHandler.Get)

	thermometersHandler := &ThermometersHandler{
		DevicePath: devicePath,
	}

	app.Register(http.MethodGet, version, thermometersPath, thermometersHandler.GetAll)

	uiHander := &UIHandler{
		FileReader: uiFiles,
	}

	app.Register(http.MethodGet, "ui", "/*filepath", uiHander.Get)

	app.Register(http.MethodOptions, "", "/", optionsHandler)

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

// TODO: re-visit this an cors.
func optionsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	return nil
}
