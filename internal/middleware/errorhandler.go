package middleware

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/pkg/errors"
)

type ErrorHandler struct {
	logger log.Logger
}

func NewErrorHandler(logger log.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (e *ErrorHandler) HandleError(next web.Handler) web.Handler {

	// Create the handler that will be attached in the middleware chain.
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		v := ctx.Value(web.CtxValuesKey).(*web.CtxValues)

		// In the event of a panic, we want to capture it here so we can send an
		// error down the stack.
		defer func() {
			if r := recover(); r != nil {

				// Indicate this request had an error.
				v.HasError = true

				// Log the panic.
				e.logger.Printf("ERROR : Panic Caught : %s\n", r)

				// Respond with the error.
				web.Error(ctx, w, errors.New("unhandled"))

				// Print out the stack.
				e.logger.Printf("ERROR : Stacktrace\n%s\n", debug.Stack())
			}
		}()

		if err := next(ctx, w, r); err != nil {

			// Indicate this request had an error.
			v.HasError = true

			// What is the root error.
			err = errors.Cause(err)

			// Log the error.
			e.logger.Printf("ERROR : %v\n", err)

			// Respond with the error.
			web.Error(ctx, w, err)

			return nil
		}

		return nil
	}

	return h
}
