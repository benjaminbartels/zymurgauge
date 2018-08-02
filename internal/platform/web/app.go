package web

import (
	"context"
	"net/http"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"

	"github.com/julienschmidt/httprouter"
)

type ctxKey int

// CtxValuesKey is the key used to save and retrieve CtxValues from the context
const CtxValuesKey ctxKey = 1

// CtxValues are context values specific to the App
type CtxValues struct {
	StartTime  time.Time
	StatusCode int
}

// Handler extends the http.HandlerFunc buy adding context param and an error
type Handler func(context.Context, http.ResponseWriter, *http.Request, Params) error

// Params implementors contain the paramaters that are extracted from an http.request
// This allows handlers to get access to the params from httprouter without it
type Params interface {
	ByName(name string) string
}

// App represents a web application that hosts http handlers that are wrapped in the provided middlewares
type App struct {
	*httprouter.Router
	logger      log.Logger
	middlewares []MiddlewareFunc
}

// NewApp creates a new App
func NewApp(logger log.Logger, mw ...MiddlewareFunc) *App {
	a := &App{
		Router:      httprouter.New(),
		logger:      logger,
		middlewares: mw,
	}
	a.registerNotFound()
	return a
}

// Register mounts the provided handler to the provided http method and path combination
func (a *App) Register(method, path string, handler Handler) {

	// Wrap handler with middlewares
	handler = wrap(handler, a.middlewares)

	// Handler function that adds the app specific values to the request context, then calls the wrapped handler
	h := func(w http.ResponseWriter, r *http.Request) {

		// Add app specific values to the request context
		ctx := context.WithValue(r.Context(), CtxValuesKey, &CtxValues{StartTime: time.Now()})

		// Get params from context
		params := httprouter.ParamsFromContext(ctx)

		// Calls the wrapped handler
		if err := handler(ctx, w, r, params); err != nil {
			if respErr := Error(ctx, w, err); err != nil {
				// ToDo: if error processing Error
				a.logger.Print(respErr)
			}
		}
	}

	// Mount the handler to the method and path combination
	a.Router.HandlerFunc(method, path, h)
}

// registerNotFound mounts a handler that is used when a path cannot be found
func (a *App) registerNotFound() {

	// Create the handler and wrap with middlewares
	f := func(ctx context.Context, w http.ResponseWriter, r *http.Request, p Params) error {
		return Error(ctx, w, ErrNotFound)
	}
	handler := wrap(f, a.middlewares)

	// Handler function that adds the app specific values to the request context, then calls the wrapped handler
	h := func(w http.ResponseWriter, r *http.Request) {

		// Add app specific values to the request context
		ctx := context.WithValue(r.Context(), CtxValuesKey, &CtxValues{StartTime: time.Now()})

		// Get params from context
		params := httprouter.ParamsFromContext(ctx)

		// Calls the wrapped handler
		if err := handler(ctx, w, r, params); err != nil {
			if respErr := Error(ctx, w, err); err != nil {
				// ToDo: if error processing Error
				a.logger.Print(respErr)
			}
		}
	}

	a.Router.NotFound = http.HandlerFunc(h)
}
