package web

import (
	"net/http"
	"strings"
)

type App struct {
	api *API
	ui  http.Handler
}

// NewApp creates a new App
func NewApp(api *API, uiFS http.FileSystem) *App {
	a := &App{
		api: api,
		ui:  http.FileServer(uiFS),
	}

	return a
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.StripPrefix("/api/", a.api).ServeHTTP(w, r)
	} else {
		http.StripPrefix("/", a.ui).ServeHTTP(w, r)
	}
}
