package app

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/rs/xid"
)

// ContextKey is how request values or stored/retrieved.
const ContextKey CtxKey = 1

const traceIDHeader = "X-Trace-ID"

// App is the http handler for the application which include the API and UI
type App struct {
	api http.Handler
	ui  http.Handler
}

// New creates a new App
func New(routes []Route, uiFS http.FileSystem, logger log.Logger) *App {
	return &App{
		api: &API{Routes: routes, Logger: logger},
		ui:  http.FileServer(uiFS),
	}
}

// Handler returns a http.Handler for the that is wrapped with the middlewares
func (a *App) Handler(middlewares ...Middleware) http.Handler {

	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/") {
			http.StripPrefix("/api/v1/", a.api).ServeHTTP(w, r)
		} else {
			http.StripPrefix("/", a.ui).ServeHTTP(w, r)
		}
	})

	for i := len(middlewares) - 1; i >= 0; i-- {
		if middlewares[i] != nil {
			mw := middlewares[i]
			handler = mw(handler)
		}
	}

	return addRequestInfo(handler)
}

// CtxKey represents the type of value for the context key
type CtxKey int

// RequestState represent state for each request.
type RequestState struct {
	TraceID    string
	StatusCode int
	Now        time.Time
}

// Middleware wraps a handler to remove boilerplate or other concerns not direct to any given Handler.
type Middleware func(http.Handler) http.Handler

func addRequestInfo(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := RequestState{
			TraceID: xid.New().String(),
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), ContextKey, &s)
		w.Header().Set(traceIDHeader, s.TraceID)
		f.ServeHTTP(w, r.WithContext(ctx))
	})
}
