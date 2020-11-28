package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// CtxValuesKey is the key used to save and retrieve CtxValues from the context.
const CtxValuesKey ctxKey = 1

type ctxKey int

// CtxValues are context values specific to the App.
type CtxValues struct {
	Now        time.Time
	StatusCode int
	HasError   bool
}

type App struct {
	router     *httprouter.Router
	middleware []MiddlewareFunc
}

// Handler extends the http.HandlerFunc buy adding context param and an error.
type Handler func(context.Context, http.ResponseWriter, *http.Request, httprouter.Params) error

func NewApp(middleware ...MiddlewareFunc) *App {
	return &App{
		router:     httprouter.New(),
		middleware: middleware,
	}
}

func (a *App) Handle(method string, path string, handler Handler, wrapWithMiddlewares bool) {
	// Wrap handler with middlewares
	if wrapWithMiddlewares {
		handler = wrap(handler, a.middleware)
	}
	// Handler function that adds the app specific values to the request context, then calls the wrapped handler
	h := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		v := &CtxValues{
			Now: time.Now(),
		}

		// Add app specific values to the request context
		ctx := context.WithValue(r.Context(), CtxValuesKey, v)

		// Calls the wrapped handler
		if err := handler(ctx, w, r, p); err != nil {
			// This is called when the error handler middleware doesn't handle the error, which is never
			fmt.Printf("ERROR : %v\n", err)
		}
	}

	// Mount the handler to the path
	a.router.Handle(method, path, h)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
