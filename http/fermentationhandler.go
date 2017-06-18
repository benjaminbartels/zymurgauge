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

// FermentationHandler is the HTTP resource handler for Fermentations.
type FermentationHandler struct {
	*httprouter.Router
	logger              *logrus.Logger
	FermentationService zymurgauge.FermentationService
}

// NewFermentationHandler instantiates a new FermentationHandler
func NewFermentationHandler(l *logrus.Logger) *FermentationHandler {
	h := &FermentationHandler{
		Router: httprouter.New(),
		logger: l,
	}
	h.POST("/api/fermentations", h.handlePost)
	h.GET("/api/fermentations/:id", h.handleGet)
	// ToDo: Move to its own service
	//h.PATCH("/api/fermentations/:id", h.handlePatch)
	return h
}

// handleGet handles the HTTP GET request for the Fermentation resource
func (h *FermentationHandler) handleGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	i := ps.ByName("id")

	id, err := strconv.ParseUint(i, 10, 64)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, h.logger)
	}

	f, err := h.FermentationService.Get(id)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError, h.logger)
	} else if f == nil {
		notFound(w, h.logger)
	} else {
		encode(w, &getFermentationResponse{Fermentation: f}, h.logger)
	}
}

// handlePost handles the HTTP POST request for the Fermentation resource
func (h *FermentationHandler) handlePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req postFermentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, zymurgauge.ErrInvalidJSON, http.StatusBadRequest, h.logger)
		return
	}
	f := req.Fermentation
	f.ModTime = time.Time{}

	switch err := h.FermentationService.Save(f); err {
	case nil:
		encode(w, &postFermentationResponse{Fermentation: f}, h.logger)
	case zymurgauge.ErrFermentationRequired:
		handleError(w, err, http.StatusBadRequest, h.logger)
	default:
		handleError(w, err, http.StatusInternalServerError, h.logger)
	}
}

// ToDo: Move to its own service
// handlePatch handles the HTTP PATCH request for the Fermentation resource
// func (h *FermentationHandler) handlePatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

// 	var req patchFermentationRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		handleError(w, zymurgauge.ErrInvalidJSON, http.StatusBadRequest)
// 		return
// 	}

// 	switch err := h.FermentationService.LogEvent(req.FermentationID, req.Event); err {
// 	case nil:
// 		encode(w, &patchFermentationResponse{})
// 	case zymurgauge.ErrNotFound:
// 		handleError(w, err, http.StatusNotFound)
// 	default:
// 		handleError(w, err, http.StatusInternalServerError)
// 	}
// }

type getFermentationResponse struct {
	Fermentation *zymurgauge.Fermentation `json:"fermentation,omitempty"`
	Err          string                   `json:"err,omitempty"`
}

type postFermentationRequest struct {
	Fermentation *zymurgauge.Fermentation `json:"fermentation,omitempty"`
}

type postFermentationResponse struct {
	Fermentation *zymurgauge.Fermentation `json:"fermentation,omitempty"`
	Err          string                   `json:"err,omitempty"`
}

type patchFermentationRequest struct {
	FermentationID uint64                       `json:"fermentationID"`
	Event          zymurgauge.FermentationEvent `json:"event"`
}

type patchFermentationResponse struct {
	Err string `json:"err,omitempty"`
}
