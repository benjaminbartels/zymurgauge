package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// RequestLogger provides a Middleware that logs out request details.
func RequestLogger(logger *logrus.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
			v, err := web.GetContextValues(ctx)
			if err != nil {
				return errors.Wrap(web.NewShutdownError(err.Error()), "could not get context values")
			}

			err = handler(ctx, w, r, p)

			logger.Infof("(%d) : %s %s -> %s (%s)", v.StatusCode, r.Method, v.Path, r.RemoteAddr, time.Since(v.Now))

			return err
		}

		return h
	}

	return m
}
