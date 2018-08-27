package middleware

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/pkg/errors"
)

// ToDo: Make ErrorHandler a HandlerFunc itself

// ErrorHandler provides a MiddlewareFunc that handles errors from handlers
type ErrorHandler struct {
	logger log.Logger
}

// NewErrorHandler creates a new ErrorHandler
func NewErrorHandler(logger log.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

// HandleError is a MiddlewareFunc that handles errors from handlers
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
				e.logger.Printf("Error : Panic Caught : %s\n", r)

				// Print out the stack.
				e.logger.Printf("Error : Stacktrace\n%s\n", debug.Stack())

				// Respond with the error.
				if err := web.Error(ctx, w, errors.New("unhandled")); err != nil {
					e.logger.Printf("Error : %s", errors.Wrap(err, "Could not send error to client")) // ToDo: Check this
				}
			}
		}()

		if err := next(ctx, w, r); err != nil {

			// Indicate this request had an error.
			v.HasError = true

			// Log the error.
			e.logger.Printf("Error : %v\n", err)

			// Respond with the error.
			err = web.Error(ctx, w, err) // ToDo: Make sure that returning err to caller is ok

			return err
		}

		return nil
	}

	return h
}
