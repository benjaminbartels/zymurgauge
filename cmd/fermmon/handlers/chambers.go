package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
	"github.com/julienschmidt/httprouter"
)

type Chambers struct {
	repo *storage.ChamberRepo
}

func (h *Chambers) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	chambers, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, chambers, http.StatusOK)
}

func (h *Chambers) Get(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	chamber, err := h.repo.Get(id)
	if err != nil {
		fmt.Println("Could not get Chambers:", err)
		os.Exit(1)
	}

	if chamber == nil {
		return web.ErrNotFound
	}

	return web.Respond(ctx, w, chamber, http.StatusOK)
}

func (h *Chambers) Save(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	chamber, err := parseChamber(r)
	if err != nil {
		return err
	}

	if err = h.repo.Save(&chamber); err != nil {
		return err
	}

	return web.Respond(ctx, w, chamber, http.StatusOK)
}

func (h *Chambers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	c, err := h.repo.Get(id)
	if err != nil {
		fmt.Println("Could not get Chambers:", err)
		os.Exit(1)
	}

	if c == nil {
		return web.ErrNotFound
	}

	if err := h.repo.Delete(id); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}

func parseChamber(r *http.Request) (storage.Chamber, error) {
	var chamber storage.Chamber
	err := json.NewDecoder(r.Body).Decode(&chamber)

	return chamber, err
}
