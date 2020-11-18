package web

import (
	"context"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
)

type ctxKey int

// CtxValuesKey is the key used to save and retrieve CtxValues from the context.
const CtxValuesKey ctxKey = 1

// CtxValues are context values specific to the App.
type CtxValues struct {
	StartTime    time.Time
	OriginalPath string
	StatusCode   int
	HasError     bool
}

// Handler extends the http.HandlerFunc buy adding context param and an error.
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

const (
	// GET method.
	GET string = "GET" // ToDo:  Are these needed?  Wrap html package completely...
	// POST method.
	POST string = "POST"
	// DELETE method.
	DELETE string = "DELETE"
)

// API implements http.HandlerFunc and represents a REST API.  It hosts a collection of routes (path and Handler pairs)
// that are used to route a request url path to the appropriate handler. Handlers are wrapped in the provided
// middleware.
type API struct {
	version     string
	routes      map[string]http.HandlerFunc
	logger      log.Logger
	middlewares []MiddlewareFunc
}

// NewAPI creates a new API withe the given version and middlewares.
func NewAPI(version string, logger log.Logger, mw ...MiddlewareFunc) *API {
	return &API{
		version:     version,
		routes:      make(map[string]http.HandlerFunc),
		logger:      logger,
		middlewares: mw,
	}
}

// ServeHTTP calls f(w, r).
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // ToDo: don't allow just *
	w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	handled := false

	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)

	// check version
	if head == a.version {
		// take off version
		r.URL.Path = strings.TrimPrefix(r.URL.Path, a.version+"/")

		// pop head off of path
		var head string
		head, r.URL.Path = ShiftPath(r.URL.Path)

		// find route by head
		handler, ok := a.routes[head]
		if ok {
			handler.ServeHTTP(w, r)

			handled = true
		}
	}

	if !handled {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Register mounts the provided handler to the provided path creating a route.
func (a *API) Register(path string, handler Handler, wrapWithMiddlewares bool) {
	// Wrap handler with middlewares
	if wrapWithMiddlewares {
		handler = wrap(handler, a.middlewares)
	}

	// Handler function that adds the app specific values to the request context, then calls the wrapped handler
	h := func(w http.ResponseWriter, r *http.Request) {
		// Calls the wrapped handler
		if err := handler(r.Context(), w, r); err != nil {
			// This is called when the error handler middleware doesn't handle the error, which is never
			a.logger.Printf("ERROR : %v\n", err)
		}
	}

	// Mount the handler to the path
	a.routes[path] = h
}

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
