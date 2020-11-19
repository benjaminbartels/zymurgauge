package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Respond sends the JSON response to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, code int) error {
	// Set the status code in context
	v := ctx.Value(CtxValuesKey).(*CtxValues)
	v.StatusCode = code

	// No Content
	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(code)

		return nil
	}

	// Marshal into a JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type, write code and write to ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonData)

	return err
}

var (
	// ErrNotFound is returned when an entity is not found.
	ErrNotFound = errors.New("not found")
	// ErrInternal is returned when an internal error has occurred.
	ErrInternal = errors.New("internal error")
	// ErrBadRequest is returned when the request is invalid.
	ErrBadRequest = errors.New("bad request")
	// ErrMethodNotAllowed is returned when the request method (GET, POST, etc.) is not allowed.
	ErrMethodNotAllowed = errors.New("method not allowed")
	// ErrUnauthorized is returned when the request is not authorized.
	ErrUnauthorized = errors.New("unauthorized")
)

// errorResponse is the response sent to the client in the event of a error.
type errorResponse struct {
	Err string `json:"error,omitempty"`
}

// Error converts application error to http error code then passes it RespondError.
func Error(ctx context.Context, w http.ResponseWriter, err error) error {
	var code int

	//nolint:errorlint
	switch errors.Cause(err) {
	case ErrNotFound:
		code = http.StatusNotFound
	case ErrBadRequest: // TODO: what was bad?
		code = http.StatusBadRequest
	case ErrMethodNotAllowed:
		code = http.StatusMethodNotAllowed
	case ErrUnauthorized:
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}

	return Respond(ctx, w, errorResponse{Err: err.Error()}, code)
}
