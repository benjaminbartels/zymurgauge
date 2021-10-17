package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// RequestError represents an error with an HTTP status code.
type RequestError struct {
	ErrMessage string
	Status     int
}

// NewRequestError creates a new RequestError with the provided error and HTTP status code.
func NewRequestError(errMessage string, status int) error {
	return &RequestError{errMessage, status}
}

// Error implements the error interface.
func (r *RequestError) Error() string {
	return r.ErrMessage
}

// Respond sends the JSON response to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	if err := SetStatusCode(ctx, statusCode); err != nil {
		return errors.Wrap(err, "could not set status code in context")
	}

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)

		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "could not marshal data")
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return errors.Wrap(err, "could not write response")
	}

	return nil
}
