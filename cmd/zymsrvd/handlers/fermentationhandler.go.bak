package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
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

// GetAll handles a GET request for all Fermentations
func (h *FermentationHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) error {
	if fermentations, err := h.repo.GetAll(); err != nil {
		return err
	} else {
		web.Respond(ctx, w, fermentations, http.StatusOK)
	}
	return nil
}

// GetOne handles a GET request for a specific Fermentation whose ID matched the provided ID
func (h *FermentationHandler) GetOne(ctx context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) error {
	id := p["id"]

	fermentationID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err //ToDo: error InvalidID
	}

	if fermentation, err := h.repo.Get(fermentationID); err != nil {
		return err
	} else if fermentation == nil {
		return web.ErrNotFound
	} else {
		web.Respond(ctx, w, fermentation, http.StatusOK)
	}
	return nil
}

// Post handles the POST request to create or update a Fermentation
func (h *FermentationHandler) Post(ctx context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) error {
	fermentation, err := parseFermentation(r)
	if err != nil {
		return err
	}

	if err := h.repo.Save(&fermentation); err != nil {
		return err
	}

	web.Respond(ctx, w, fermentation, http.StatusOK)
	return nil
}

// Delete handles the DELETE request to delete a Fermentation
func (h *FermentationHandler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) error {

	id := p["id"]

	fermentationID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err //ToDo: error InvalidID
	}

	if err := h.repo.Delete(fermentationID); err != nil {
		return err
	}

	web.Respond(ctx, w, nil, http.StatusOK)
	return nil
}

// parseFermentation decodes the specified Fermentation into JSON
func parseFermentation(r *http.Request) (internal.Fermentation, error) {
	var fermentation internal.Fermentation
	err := json.NewDecoder(r.Body).Decode(&fermentation)
	return fermentation, err
}
