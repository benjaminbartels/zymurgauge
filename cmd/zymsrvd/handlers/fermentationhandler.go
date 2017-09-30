package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/platform/app"
)

// FermentationHandler is the http handler for API calls to manage Fermentations
type FermentationHandler struct {
	repo *database.FermentationRepo
}

// NewFermentationHandler instantiates a FermentationHandler
func NewFermentationHandler(repo *database.FermentationRepo) *FermentationHandler {
	return &FermentationHandler{
		repo: repo,
	}
}

// ServeHTTP calls f(w, r).
func (h *FermentationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGet(w, r)
	case "POST":
		h.handlePost(w, r)
	default:
		app.HandleError(w, app.ErrNotFound)
	}
}

func (h *FermentationHandler) handleGet(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "" {

		if id, err := strconv.ParseUint(r.URL.Path, 10, 64); err != nil {
			app.HandleError(w, app.ErrBadRequest)
		} else {
			h.handleGetOne(w, id)
		}

	} else {
		h.handleGetAll(w)
	}

}

func (h *FermentationHandler) handleGetOne(w http.ResponseWriter, id uint64) {
	if fermentation, err := h.repo.Get(id); err != nil {
		app.HandleError(w, err)
	} else if fermentation == nil {
		app.HandleError(w, app.ErrNotFound)
	} else {
		app.Encode(w, &fermentation)
	}

}

func (h *FermentationHandler) handleGetAll(w http.ResponseWriter) {
	if fermentations, err := h.repo.GetAll(); err != nil {
		app.HandleError(w, err)
	} else {
		app.Encode(w, fermentations)
	}
}

func (h *FermentationHandler) handlePost(w http.ResponseWriter, r *http.Request) {

	fermentation, err := parseFermentation(r)
	if err != nil {
		app.HandleError(w, err)
		return
	}

	if err := h.repo.Save(&fermentation); err != nil {
		app.HandleError(w, err)
	} else {
		app.Encode(w, &fermentation)
	}
}

func parseFermentation(r *http.Request) (internal.Fermentation, error) {
	var fermentation internal.Fermentation
	err := json.NewDecoder(r.Body).Decode(&fermentation)
	return fermentation, err
}
