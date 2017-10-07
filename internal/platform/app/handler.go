package app

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
)

// Handler is imbedded by implementors of http.HandlerFunc for easy access
// to encoding to the ResponseWriter
type Handler struct {
	Logger log.Logger
}

// Encode encodes the given interface onto the given http.ResponseWriter
func (h *Handler) Encode(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.HandleError(w, err)
	}
}

// HandleError encodes the given error onto the given http.ResponseWriter
func (h *Handler) HandleError(w http.ResponseWriter, err error) {
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
		h.Logger.Println(encErr)
	}
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
func (h Handler) ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
