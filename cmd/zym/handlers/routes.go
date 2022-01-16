package handlers

import (
	"context"
	"embed"
	"expvar"
	"fmt"
	"mime"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"strings"

	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/middleware"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
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
func NewAPI(chamberManager chamber.Controller, devicePath string, service brewfather.Service, uiFiles embed.FS,
	shutdown chan os.Signal, logger *logrus.Logger) http.Handler {
	chambersHandler := &ChambersHandler{
		ChamberController: chamberManager,
		Logger:            logger,
	}

	thermometersHandler := &ThermometersHandler{
		DevicePath: devicePath,
	}

	batchesHandler := &BatchesHandler{
		Service: service,
	}

	uiHander := &UIHandler{
		UIFiles: uiFiles,
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

	app.Register(http.MethodGet, "ui", "/*filepath", uiHander.handle)

	return app
}

type UIHandler struct {
	UIFiles embed.FS
}

func (h *UIHandler) handle(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var filePath string

	if strings.HasPrefix(r.URL.Path, base+"/static") {
		filePath = uiDir + strings.TrimPrefix(r.URL.Path, base)
	} else {
		filePath = uiDir + "/index.html"
	}

	b, err := h.UIFiles.ReadFile(filePath)
	if err != nil {
		return errors.Wrap(err, "could not read file")
	}

	if contentType := mime.TypeByExtension(filepath.Ext(r.URL.Path)); len(contentType) > 0 {
		w.Header().Add("Content-Type", contentType)
	}

	if _, err := w.Write(b); err != nil {
		return errors.Wrap(err, "could not write response")
	}

	return nil
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
