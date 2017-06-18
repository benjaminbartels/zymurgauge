package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/benjaminbartels/zymurgauge"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// BeerHandler is the HTTP resource handler for Beers.
type BeerHandler struct {
	*httprouter.Router
	logger      *logrus.Logger
	BeerService zymurgauge.BeerService
}

// NewBeerHandler instantiates a new BeerHandler
func NewBeerHandler(l *logrus.Logger) *BeerHandler {
	h := &BeerHandler{
		Router: httprouter.New(),
		logger: l,
	}
	h.POST("/api/beers", h.handlePost)
	h.GET("/api/beers/:id", h.handleGet)
	return h
}

// handleGet handles the HTTP GET request for the Beer resource
func (h *BeerHandler) handleGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	i := ps.ByName("id")

	id, err := strconv.ParseUint(i, 10, 64)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, h.logger)
	}

	b, err := h.BeerService.Get(id)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError, h.logger)
	} else if b == nil {
		notFound(w, h.logger)
	} else {
		encode(w, &getBeerResponse{Beer: b}, h.logger)
	}
}

// handlePost handles the HTTP POST request for the Beer resource
func (h *BeerHandler) handlePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req postBeerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, zymurgauge.ErrInvalidJSON, http.StatusBadRequest, h.logger)
		return
	}
	b := req.Beer
	b.ModTime = time.Time{}

	switch err := h.BeerService.Save(b); err {
	case nil:
		encode(w, &postBeerResponse{Beer: b}, h.logger)
	case zymurgauge.ErrBeerRequired:
		handleError(w, err, http.StatusBadRequest, h.logger)
	default:
		handleError(w, err, http.StatusInternalServerError, h.logger)
	}
}

type getBeerResponse struct {
	Beer *zymurgauge.Beer `json:"beer,omitempty"`
	Err  string           `json:"err,omitempty"`
}

type postBeerRequest struct {
	Beer *zymurgauge.Beer `json:"beer,omitempty"`
}

type postBeerResponse struct {
	Beer *zymurgauge.Beer `json:"beer,omitempty"`
	Err  string           `json:"err,omitempty"`
}
