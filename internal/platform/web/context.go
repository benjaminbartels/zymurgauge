package web

import (
	"context"
	"errors"
	"time"
)

var errMissingInContext = errors.New("value missing from context")

type ctxKey int

const key ctxKey = 1

// CtxValues represent state for each request.
type CtxValues struct {
	Path       string
	Now        time.Time
	StatusCode int
}

// GetValues returns the values from the context.
func GetContextValues(ctx context.Context) (*CtxValues, error) {
	v, ok := ctx.Value(key).(*CtxValues)
	if !ok {
		return nil, errMissingInContext
	}

	return v, nil
}

// InitContextValues initializes the CtxValues in the context with the given values and return the updated context.
func InitContextValues(ctx context.Context, v *CtxValues) context.Context {
	return context.WithValue(ctx, key, v)
}

// SetStatusCode sets the status code on the context.
func SetStatusCode(ctx context.Context, statusCode int) error {
	v, ok := ctx.Value(key).(*CtxValues)
	if !ok {
		return errMissingInContext
	}

	v.StatusCode = statusCode

	return nil
}
