package app

import "errors"

var (
	// ErrNotFound is returned when an entity is not found
	ErrNotFound = errors.New("not found")
	// ErrNotAllowed is returned when the http method is not allowed
	ErrNotAllowed = errors.New("method not allowed")
	// ErrBadRequest is returned when a bad request has occurred
	ErrBadRequest = errors.New("bad request")
	// ErrInternal is returned when an internal error has occurred
	ErrInternal = errors.New("internal error")
	// ErrInvalidJSON is returned when json request is invalid
	ErrInvalidJSON = errors.New("invalid json")
)

type errorResponse struct {
	Err string `json:"error,omitempty"`
}
