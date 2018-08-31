package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

// ToDo: Make RequestLogger a HandlerFunc itself

// RequestLogger provides a MiddlewareFunc that logs out request details
type RequestLogger struct {
	logger log.Logger
}

// NewRequestLogger creates a new RequestLogger
func NewRequestLogger(logger log.Logger) *RequestLogger {
	return &RequestLogger{logger: logger}
}

// Log is a MiddlewareFunc that logs out request details including app specific context values
func (l *RequestLogger) Log(next web.Handler) web.Handler {

	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		err := next(ctx, w, r)

		v := ctx.Value(web.CtxValuesKey).(*web.CtxValues)

		// ToDo: r.URL.Path is getting overwritten, create a value in ctx?
		l.logger.Printf("(%d) : %s %s -> %s (%s)", v.StatusCode, r.Method, r.URL.Path, r.RemoteAddr,
			time.Since(v.StartTime))

		return err

	}

	return h
}
