package web

import (
	"context"
	"net/http"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
)

// App represents a web application that includes a REST API and a hosted web based UI.
type App struct {
	api    *API
	logger log.Logger
}

// NewApp creates a new App.  The usFS is a filesystem that gets mounted into a http.FileServer.
func NewApp(api *API, logger log.Logger) *App {
	a := &App{
		api:    api,
		logger: logger,
	}
	return a
}

// ServeHTTP calls f(w, r). API calls are routed to the API. Calls at the root "/" are intended to be for the UI and are
// routed to an FileServer.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := &CtxValues{
		StartTime:    time.Now(),
		OriginalPath: r.URL.Path,
	}

	// Add app specific values to the request context
	ctx := context.WithValue(r.Context(), CtxValuesKey, v)
	_, r.URL.Path = ShiftPath(r.URL.Path)
	a.api.ServeHTTP(w, r.WithContext(ctx))
}
