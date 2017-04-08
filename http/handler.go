package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/orangesword/zymurgauge"
	"github.com/sirupsen/logrus"
)

// Handler contains the references to the service handlers
type Handler struct {
	ChamberHandler      *ChamberHandler
	FermentationHandler *FermentationHandler
	BeerHandler         *BeerHandler
}

type errorResponse struct {
	Err string `json:"err,omitempty"`
}

// ServeHTTP delegates requests to the service handlers.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, "/api/chambers") {
		h.ChamberHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/fermentations") {
		h.FermentationHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/beers") {
		h.BeerHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// encode encodes the given interface to the given response writer
func encode(w http.ResponseWriter, v interface{}, logger *logrus.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		handleError(w, err, http.StatusInternalServerError, logger)
	}

}

// handleError encodes the given error to the given response writer
func handleError(w http.ResponseWriter, err error, code int, logger *logrus.Logger) {

	logger.Infof("http error: %s (code=%d)", err, code)

	if code == http.StatusInternalServerError {
		err = zymurgauge.ErrInternal
	}

	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
	if err != nil {
		logger.Error(err)
	}
}

// notFound encodes and Not Found status to the header and a empty JSON object to the response
func notFound(w http.ResponseWriter, logger *logrus.Logger) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte(`{}` + "\n"))
	if err != nil {
		logger.Error(err)
	}
}
