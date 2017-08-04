package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ErrNotFound    = Error("not found")
	ErrNotAllowed  = Error("method not allowed")
	ErrBadRequest  = Error("bad request")
	ErrInternal    = Error("internal error")
	ErrInvalidJSON = Error("invalid json")
)

type Error string

func (e Error) Error() string { return string(e) }

// ToDo: move to platform?
func Encode(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		HandleError(w, err)
	}
}

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
