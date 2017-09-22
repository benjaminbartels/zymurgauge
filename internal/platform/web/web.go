package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

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

// Encode encodes the given interface onto the given http.ResponseWriter
func Encode(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		HandleError(w, err)
	}
}

// HandleError encodes the given error onto the given http.ResponseWriter
func HandleError(w http.ResponseWriter, err error) {
	var code int

	switch err {
	case ErrNotFound:
		code = http.StatusNotFound
	case ErrNotAllowed:
		code = http.StatusMethodNotAllowed
	case ErrBadRequest:
		code = http.StatusBadRequest
	case ErrInvalidJSON:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}

	w.WriteHeader(code)
	encErr := json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
	if encErr != nil {
		fmt.Println(encErr)
	}
}

type errorResponse struct {
	Err string `json:"error,omitempty"`
}
