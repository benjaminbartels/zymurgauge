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

// Handle handles the incoming http request
func (h *BeerHandler) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case web.GET:
		return h.get(ctx, w, r)
	case web.POST:
		return h.post(ctx, w, r)
	case web.DELETE:
		return h.delete(r)
	default:
		return web.ErrMethodNotAllowed
	}
}

func (h *BeerHandler) get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var head string
	head, r.URL.Path = web.ShiftPath(r.URL.Path)
	if head == "" {
		return h.getAll(ctx, w)
	}
	id, err := strconv.ParseUint(head, 10, 64)
	if err != nil {
		return err
	}
	return h.getOne(ctx, w, id)
}

func (h *BeerHandler) getAll(ctx context.Context, w http.ResponseWriter) error {
	beers, err := h.repo.GetAll()
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, beers, http.StatusOK)
}

func (h *BeerHandler) getOne(ctx context.Context, w http.ResponseWriter, id uint64) error {
	if beer, err := h.repo.Get(id); err != nil {
		return err
	} else if beer == nil {
		return web.ErrNotFound
	} else {
		return web.Respond(ctx, w, beer, http.StatusOK)
	}
}

func (h *BeerHandler) post(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	beer, err := parseBeer(r)
	if err != nil {
		return err
	}
	err = h.repo.Save(&beer)
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, beer, http.StatusOK)
}

func (h *BeerHandler) delete(r *http.Request) error {
	if r.URL.Path == "" {
		return web.ErrBadRequest
	}
	id, err := strconv.ParseUint(r.URL.Path, 10, 64)
	if err != nil {
		return err
	}
	if err := h.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

func parseBeer(r *http.Request) (internal.Beer, error) {
	var beer internal.Beer
	err := json.NewDecoder(r.Body).Decode(&beer)
	return beer, err
}
