package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Respond sends the JSON response to the client
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
		if respErr := respondError(ctx, w, err, http.StatusInternalServerError); respErr != nil {
			return respErr
		}

		return err
	}

	// Set the content type, write code and write to ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonData)
	return err
}

var (
	// ErrNotFound is returned when an entity is not found
	ErrNotFound = errors.New("not found")
	// ErrInternal is returned when an internal error has occurred
	ErrInternal = errors.New("internal error")
)

// errorResponse is the response sent to the client in the event of a error
type errorResponse struct {
	Err string `json:"error,omitempty"`
}

// Error converts application error to http error code then passes it RespondError
func Error(ctx context.Context, w http.ResponseWriter, err error) error {
	switch errors.Cause(err) {
	case ErrNotFound:
		return respondError(ctx, w, err, http.StatusNotFound)
	}
	return respondError(ctx, w, err, http.StatusInternalServerError)
}

// respondError sends the JSON error response to the client
func respondError(ctx context.Context, w http.ResponseWriter, err error, code int) error {
	return Respond(ctx, w, errorResponse{Err: err.Error()}, code)
}
