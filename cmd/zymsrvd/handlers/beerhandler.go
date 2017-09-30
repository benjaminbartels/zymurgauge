package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/platform/app"
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

// ServeHTTP calls f(w, r).
func (h *BeerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGet(w, r)
	case "POST":
		h.handlePost(w, r)
	default:
		app.HandleError(w, app.ErrNotFound)
	}
}

func (h *BeerHandler) handleGet(w http.ResponseWriter, r *http.Request) {
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

func (h *BeerHandler) handleGetOne(w http.ResponseWriter, id uint64) {
	if beer, err := h.repo.Get(id); err != nil {
		app.HandleError(w, err)
	} else if beer == nil {
		app.HandleError(w, app.ErrNotFound)
	} else {
		app.Encode(w, &beer)
	}

}

func (h *BeerHandler) handleGetAll(w http.ResponseWriter) {
	if beers, err := h.repo.GetAll(); err != nil {
		app.HandleError(w, err)
	} else {
		app.Encode(w, beers)
	}
}

func (h *BeerHandler) handlePost(w http.ResponseWriter, r *http.Request) {

	beer, err := parseBeer(r)
	if err != nil {
		app.HandleError(w, err)
		return
	}

	if err := h.repo.Save(&beer); err != nil {
		app.HandleError(w, err)
	} else {
		app.Encode(w, &beer)
	}
}

func parseBeer(r *http.Request) (internal.Beer, error) {
	var beer internal.Beer
	err := json.NewDecoder(r.Body).Decode(&beer)
	return beer, err
}
