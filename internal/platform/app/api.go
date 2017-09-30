package app

import (
	"fmt"
	"net/http"
	"strings"
)

// API is the http handler for call to the API
type API struct {
	Routes []Route
}

// ServeHTTP calls f(w, r)
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	handled := false
	for _, route := range a.Routes {
		fmt.Println("Checking", route.Path)
		if strings.HasPrefix(r.URL.Path, route.Path+"/") {
			fmt.Println("has /")
			handled = true
			http.StripPrefix(route.Path+"/", route.Handler).ServeHTTP(w, r)
		} else if r.URL.Path == route.Path {
			fmt.Println("is equal to ")
			handled = true
			http.StripPrefix(route.Path, route.Handler).ServeHTTP(w, r)
		}
	}

	if !handled {
		HandleError(w, ErrNotFound)
	}
}

// Route associates a path to a http handler
type Route struct {
	Path    string
	Handler http.Handler
}
