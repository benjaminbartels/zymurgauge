package web

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	uiDir = "build"
)

type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type App struct {
	api          *API
	uiFileReader FileReader
	mux          *http.ServeMux
	logger       *logrus.Logger
}

func NewApp(api *API, uiFileReader FileReader, logger *logrus.Logger) *App {
	app := &App{
		api:          api,
		uiFileReader: uiFileReader,
		logger:       logger,
	}

	app.mux = http.NewServeMux()
	app.mux.Handle("/api/", api)
	app.mux.Handle("/", http.HandlerFunc(app.ui))

	return app
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *App) ui(w http.ResponseWriter, r *http.Request) {
	var filePath string

	switch {
	case strings.HasPrefix(r.URL.Path, "/static"):
		filePath = uiDir + r.URL.Path
	case strings.HasSuffix(r.URL.Path, ".png"):
		filePath = uiDir + r.URL.Path
	case strings.HasSuffix(r.URL.Path, ".ico"):
		filePath = uiDir + r.URL.Path
	default:
		filePath = uiDir + "/index.html"
	}

	b, err := a.uiFileReader.ReadFile(filePath)
	if err != nil {
		a.logger.WithError(err).Error("Could not read file")
	}

	if contentType := mime.TypeByExtension(filepath.Ext(r.URL.Path)); len(contentType) > 0 {
		h := w.Header()
		h.Add("Content-Type", contentType)
	}

	if _, err := w.Write(b); err != nil {
		a.logger.WithError(err).Error("Could not write response")
	}
}
