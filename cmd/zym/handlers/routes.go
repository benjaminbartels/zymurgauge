package handlers

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/controller"
	"github.com/benjaminbartels/zymurgauge/internal/batch"
	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/middleware"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/sirupsen/logrus"
)

const (
	chambersPath     = "/chambers"
	thermometersPath = "/thermometers"
	batchesPath      = "/batches"
	version          = "v1"
)

// NewAPI return a web.App with configured routes and handlers.
func NewAPI(chamberController controller.ChamberController, thermometerRepo device.ThermometerRepo,
	batchRepo batch.Repo, shutdown chan os.Signal, logger *logrus.Logger) http.Handler {
	chambersHandler := &ChambersHandler{
		ChamberController: chamberController,
		Logger:            logger,
	}

	thermometersHandler := &ThermometersHandler{
		Repo: thermometerRepo,
	}

	batchesHandler := &BatchesHandler{
		Repo: batchRepo,
	}

	app := web.NewApp(shutdown, middleware.RequestLogger(logger), middleware.Errors(logger))

	app.Register(http.MethodGet, version, chambersPath, chambersHandler.GetAll)
	app.Register(http.MethodGet, version, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Get)
	app.Register(http.MethodPost, version, chambersPath, chambersHandler.Save)
	app.Register(http.MethodDelete, version, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Delete)
	app.Register(http.MethodPost, version, fmt.Sprintf("%s/:id/start", chambersPath), chambersHandler.Start)
	app.Register(http.MethodPost, version, fmt.Sprintf("%s/:id/stop", chambersPath), chambersHandler.Stop)

	app.Register(http.MethodGet, version, thermometersPath, thermometersHandler.GetAll)

	app.Register(http.MethodGet, version, batchesPath, batchesHandler.GetAll)
	app.Register(http.MethodGet, version, fmt.Sprintf("%s/:id", batchesPath), batchesHandler.Get)

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
