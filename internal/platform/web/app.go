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

// App represents a web application that hosts a REST API.
type App struct {
	router      *httprouter.Router
	shutdown    chan os.Signal
	middlewares []Middleware
}

// NewApp creates an App that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, middlewares ...Middleware) *App {
	router := httprouter.New()

	return &App{
		router:      router,
		shutdown:    shutdown,
		middlewares: middlewares,
	}
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// Register mounts the provided handler to the provided path creating a route.
func (a *App) Register(method, group, path string, handler Handler) {
	handler = wrap(a.middlewares, handler)

	h := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		v := CtxValues{
			Path: r.URL.Path,
			Now:  time.Now(),
		}

		ctx := InitContextValues(r.Context(), &v)

		if err := handler(ctx, w, r, p); err != nil {
			// TODO: log here?
			a.SignalShutdown() // TODO: is this necessary?

			return
		}
	}

	if group != "" {
		path = "/" + group + path
	}

	a.router.Handle(method, path, h)
}

// SignalShutdown is used to gracefully shutdown the app.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
