package web

import (
	"context"
	"net/http"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
)

// App represents a web application that includes a REST API and a hosted web based UI
type App struct {
	api    *API
	fs     http.FileSystem
	ui     http.Handler
	logger log.Logger
}

// NewApp creates a new App.  The usFS is a filesystem that gets mounted into a http.FileServer
func NewApp(api *API, uiFS http.FileSystem, logger log.Logger) *App {
	a := &App{
		api:    api,
		fs:     uiFS,
		ui:     http.FileServer(uiFS),
		logger: logger,
	}

	return a
}

// ServeHTTP calls f(w, r). API calls are routed to the API. Calls at the root "/" are intended to be for the UI and are
// routed to an FileServer
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	v := &CtxValues{
		StartTime:    time.Now(),
		OriginalPath: r.URL.Path,
	}

	// Add app specific values to the request context
	ctx := context.WithValue(r.Context(), CtxValuesKey, v)

	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)

	if head == "api" {
		// Handle calls to the API
		a.api.ServeHTTP(w, r.WithContext(ctx))
	} else if head == "static" {
		// Handle calls to static ui files
		http.StripPrefix("/", a.ui).ServeHTTP(w, r) // ToDo: wrap handlers here to use middleware
	} else {
		// Everything else gets routed to index.html
		// We can't tell the http.FileServer to serve a specific file, so we do it http.ServeContent
		f, err := a.fs.Open("/index.html")
		if err != nil {
			a.logger.Println(err)
		} else {
			http.ServeContent(w, r, "index.html", time.Time{}, f) // ToDo: What set modtime to?
		}
	}
}
