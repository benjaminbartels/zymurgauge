package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

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
func New(routes []Route, uiFS http.FileSystem) *App {
	return &App{
		api: &API{Routes: routes},
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

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

// Encode encodes the given interface onto the given http.ResponseWriter
func Encode(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		HandleError(w, err)
	}
}

// HandleError encodes the given error onto the given http.ResponseWriter
func HandleError(w http.ResponseWriter, err error) {
	var code int

	switch err {
	case ErrNotFound:
		code = http.StatusNotFound
	case ErrNotAllowed:
		code = http.StatusMethodNotAllowed
	case ErrBadRequest:
		code = http.StatusBadRequest
	case ErrInvalidJSON:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}

	w.WriteHeader(code)
	encErr := json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
	if encErr != nil {
		fmt.Println(encErr) // This should never happen...
	}
}

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
