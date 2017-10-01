package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/app"
)

// RequestLogger writes information about the request to the logs in
// the format: TraceID : (200) GET /foo -> IP ADDR (latency)
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := r.Context().Value(app.ContextKey).(*app.RequestState)
		next.ServeHTTP(w, r)
		log.Printf("%s : (%d) : %s %s -> %s (%s)",
			s.TraceID,
			s.StatusCode,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(s.Now),
		)
	})
}
