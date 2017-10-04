package middleware

import (
	"net/http"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/app"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
)

// RequestLogger is middle ware used to log details about the http request
type RequestLogger struct {
	logger log.Logger
}

// NewRequestLogger create a new RequestLogger using the provided logger
func NewRequestLogger(logger log.Logger) *RequestLogger {
	return &RequestLogger{logger}
}

// Handler writes information about the request to the logs in
// the format: TraceID : (200) GET /foo -> IP ADDR (latency)
func (l *RequestLogger) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := r.Context().Value(app.ContextKey).(*app.RequestState)
		next.ServeHTTP(w, r)
		l.logger.Printf("%s : (%d) : %s %s -> %s (%s)",
			s.TraceID,
			s.StatusCode,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(s.Now),
		)
	})
}
