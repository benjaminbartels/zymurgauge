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

func (h *BeerHandler) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case web.GET:
		return h.handleGet(ctx, w, r)
	case web.POST:
		return h.handlePost(ctx,w, r)
	case web.DELETE:
		return h.handleDelete(ctx,w, r)
	default:
		return web.ErrMethodNotAllowed
	}
}

func (h *BeerHandler) handleGet(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var head string
	head, r.URL.Path = web.ShiftPath(r.URL.Path)
	if head == "" {
		return h.handleGetAll(ctx, w)
	} else {
		if id, err := strconv.ParseUint(head, 10, 64); err != nil {
			return err
		} else {
			return h.handleGetOne(ctx, w, id)
		}
	}
}

func (h *BeerHandler) handleGetAll(ctx context.Context, w http.ResponseWriter) error {
	if beers, err := h.repo.GetAll(); err != nil {
		return err
	} else {
		return web.Respond(ctx, w, beers, http.StatusOK)
	}
}

func (h *BeerHandler) handleGetOne(ctx context.Context, w http.ResponseWriter, id uint64) error {
	if beer, err := h.repo.Get(id); err != nil {
		return err
	} else if beer == nil {
		return web.ErrNotFound
	} else {
		return web.Respond(ctx, w, beer, http.StatusOK)
	}
}

func (h *BeerHandler) handlePost(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	beer, err := parseBeer(r)
	if err != nil {
		return err
	}
	if err := h.repo.Save(&beer); err != nil {
		return err
	} else {
		return web.Respond(ctx, w, beer, http.StatusOK)
	}
}

func (h *BeerHandler) handleDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "" {
		return web.ErrBadRequest
	}
	if id, err := strconv.ParseUint(r.URL.Path, 10, 64); err != nil {
		return err
	} else {
		if err := h.repo.Delete(id); err != nil {
			return err
		}
	}
	return nil
}

func parseBeer(r *http.Request) (internal.Beer, error) {
	var beer internal.Beer
	err := json.NewDecoder(r.Body).Decode(&beer)
	return beer, err
}
