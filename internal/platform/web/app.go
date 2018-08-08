package web

import (
	"net/http"
	"strings"
)

// App represents a web application that includes a REST API and a hosted web based UI
type App struct {
	api *API
	ui  http.Handler
}

// NewApp creates a new App.  The usFS is a filesystem that gets mounted into a http.FileServer
func NewApp(api *API, uiFS http.FileSystem) *App {
	a := &App{
		api: api,
		ui:  http.FileServer(uiFS),
	}

	return a
}

// ServeHTTP calls f(w, r). API calls are routed to the API. Calls at the root "/" are intended to be for the UI and are
// routed to an FileServer
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.StripPrefix("/api/", a.api).ServeHTTP(w, r)
	} else {
		http.StripPrefix("/", a.ui).ServeHTTP(w, r)
	}
}
