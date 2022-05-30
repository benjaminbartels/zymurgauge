// Package web contains a small web framework extension.
package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Handler extends the http.HandlerFunc buy adding a context, params and an error to return.
type Handler func(context.Context, http.ResponseWriter, *http.Request, httprouter.Params) error

// API represents a web application that hosts a REST API.
type API struct {
	router      *httprouter.Router
	shutdown    chan os.Signal
	middlewares []Middleware
}

// NewAPI creates an API that handle a set of routes for the application.
func NewAPI(shutdown chan os.Signal, middlewares ...Middleware) *API {
	router := httprouter.New()

	return &API{
		router:      router,
		shutdown:    shutdown,
		middlewares: middlewares,
	}
}

// ServeHTTP implements the http.Handler interface.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// Register mounts the provided handler to the provided path creating a route.
func (a *API) Register(method, group, path string, handler Handler, middlewares ...Middleware) {
	handler = wrap(middlewares, handler)
	handler = wrap(a.middlewares, handler)

	h := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		v := CtxValues{
			Path: r.URL.Path,
			Now:  time.Now(),
		}

		ctx := InitContextValues(r.Context(), &v)

		if err := handler(ctx, w, r, p); err != nil {
			a.SignalShutdown()

			return
		}
	}

	if group != "" {
		path = "/api/" + group + path
	}

	a.router.Handle(method, path, h)
}

// SignalShutdown is used to gracefully shutdown the API.
func (a *API) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
