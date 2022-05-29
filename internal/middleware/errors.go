package middleware

import (
	"context"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Errors handles errors coming out of handlers.
func Errors(logger *logrus.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
			if err := handler(ctx, w, r, p); err != nil {
				errResponse := struct {
					Err string `json:"error"`
				}{}

				var status int

				var requestError *web.RequestError

				if errors.As(err, &requestError) {
					// handler has set this error as a request error.
					errResponse.Err = requestError.Error()
					status = requestError.Status
				} else {
					// error is unexpected and is set to
					logger.Error(err)
					errResponse.Err = http.StatusText(http.StatusInternalServerError)
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, errResponse, status); err != nil {
					return errors.Wrap(err, "could not send response")
				}

				if ok := web.IsShutdown(err); ok {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
