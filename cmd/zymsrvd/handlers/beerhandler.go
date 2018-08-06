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

// BeerHandler is the http handler for API calls to manage Beers
type BeerHandler struct {
	repo *database.BeerRepo
}

// NewBeerHandler instantiates a BeerHandler
func NewBeerHandler(repo *database.BeerRepo) *BeerHandler {
	return &BeerHandler{
		repo: repo,
	}
}

// GetAll handles a GET request for all Beers
func (h *BeerHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) error {
	if beers, err := h.repo.GetAll(); err != nil {
		return err
	} else {
		return web.Respond(ctx, w, beers, http.StatusOK)
	}
	return nil
}

// GetOne handles a GET request for a specific Beer whose ID matched the provided ID
func (h *BeerHandler) GetOne(ctx context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) error {
	id := p["id"]
	beerID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err //ToDo: error InvalidID
	}

	beer, err := h.repo.Get(beerID)

	if err != nil {
		return err
	} else if beer == nil {
		err = web.ErrNotFound
		return err
	}

	return web.Respond(ctx, w, beer, http.StatusOK)
	//return nil
}

// Post handles the POST request to create or update a Beer
func (h *BeerHandler) Post(ctx context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) error {
	beer, err := parseBeer(r)
	if err != nil {
		return err
	}

	if err := h.repo.Save(&beer); err != nil {
		return err
	}

	return web.Respond(ctx, w, beer, http.StatusOK)
	//return nil

}

// Delete handles the DELETE request to delete a Beer
func (h *BeerHandler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) error {
	id := p["id"]
	beerID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err //ToDo: error InvalidID
	}

	if err := h.repo.Delete(beerID); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
	//return nil
}

// parseBeer decodes the specified Beer into JSON
func parseBeer(r *http.Request) (internal.Beer, error) {
	var beer internal.Beer
	err := json.NewDecoder(r.Body).Decode(&beer)
	return beer, err
}
