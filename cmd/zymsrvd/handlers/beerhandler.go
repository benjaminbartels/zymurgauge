package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/platform/app"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
)

// BeerHandler is the http handler for API calls to manage Beers
type BeerHandler struct {
	app.Handler
	repo *database.BeerRepo
}

// NewBeerHandler instantiates a BeerHandler
func NewBeerHandler(repo *database.BeerRepo, logger log.Logger) *BeerHandler {
	return &BeerHandler{
		Handler: app.Handler{Logger: logger},
		repo:    repo,
	}
}

// ServeHTTP calls f(w, r).
func (h *BeerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case app.GET:
		h.handleGet(w, r)
	case app.POST:
		h.handlePost(w, r)
	default:
		h.HandleError(w, app.ErrNotFound)
	}
}

func (h *BeerHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "" {
		if id, err := strconv.ParseUint(r.URL.Path, 10, 64); err != nil {
			h.HandleError(w, app.ErrBadRequest)
		} else {
			h.handleGetOne(w, id)
		}

	} else {
		h.handleGetAll(w)
	}
}

func (h *BeerHandler) handleGetOne(w http.ResponseWriter, id uint64) {
	if beer, err := h.repo.Get(id); err != nil {
		h.HandleError(w, err)
	} else if beer == nil {
		h.HandleError(w, app.ErrNotFound)
	} else {
		h.Encode(w, &beer)
	}

}

func (h *BeerHandler) handleGetAll(w http.ResponseWriter) {
	if beers, err := h.repo.GetAll(); err != nil {
		h.HandleError(w, err)
	} else {
		h.Encode(w, beers)
	}
}

func (h *BeerHandler) handlePost(w http.ResponseWriter, r *http.Request) {

	beer, err := parseBeer(r)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	if err := h.repo.Save(&beer); err != nil {
		h.HandleError(w, err)
	} else {
		h.Encode(w, &beer)
	}
}

func parseBeer(r *http.Request) (internal.Beer, error) {
	var beer internal.Beer
	err := json.NewDecoder(r.Body).Decode(&beer)
	return beer, err
}
