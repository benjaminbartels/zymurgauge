package handlers

import (
	"net/http"
	"path"
	"strings"
)

// App is the http handler for call to the Application which include teh API and WebApp
type App struct {
	http.Handler
	API *API
	Web http.Handler
}

// ServeHTTP calls f(w, r).
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, "/api") {
		_, r.URL.Path = shiftPath(r.URL.Path)
		a.API.ServeHTTP(w, r)
	} else {
		a.Web.ServeHTTP(w, r)
	}
}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
