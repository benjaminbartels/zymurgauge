package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
)

type App struct {
	api http.Handler
	ui  http.Handler
}

func New(routes []Route, uiFS http.FileSystem) *App {
	return &App{
		api: &API{Routes: routes},
		ui:  http.FileServer(uiFS),
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/v1/") {
		http.StripPrefix("/api/v1/", a.api).ServeHTTP(w, r)
	} else {
		http.StripPrefix("/", a.ui).ServeHTTP(w, r)
	}
}

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

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
