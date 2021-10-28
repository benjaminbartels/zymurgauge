package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type ChambersHandler struct {
	Repo internal.ChamberRepo
}

func (h *ChambersHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params) error {
	chambers, err := h.Repo.GetAll()
	if err != nil {
		return errors.Wrap(err, "could not get all chambers from repository")
	}

	if err = web.Respond(ctx, w, chambers, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *ChambersHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	chamber, err := h.Repo.Get(id)
	if err != nil {
		return errors.Wrap(err, "could not get chamber from repository")
	}

	if chamber == nil {
		return web.NewRequestError(fmt.Sprintf("chamber '%s' not found", id), http.StatusNotFound)
	}

	if err = web.Respond(ctx, w, chamber, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *ChambersHandler) Save(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	chamber, err := parseChamber(r)
	if err != nil {
		return errors.Wrap(err, "could not parse chamber")
	}

	if err = h.Repo.Save(&chamber); err != nil {
		return errors.Wrap(err, "could not save chamber to repository")
	}

	if err = web.Respond(ctx, w, chamber, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *ChambersHandler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params) error {
	id := p.ByName("id")

	c, err := h.Repo.Get(id)
	if err != nil {
		return errors.Wrap(err, "could not get chamber from repository")
	}

	if c == nil {
		return web.NewRequestError(fmt.Sprintf("chamber '%s' not found", id), http.StatusNotFound)
	}

	if err := h.Repo.Delete(id); err != nil {
		return errors.Wrap(err, "could not delete chamber from repository")
	}

	if err = web.Respond(ctx, w, nil, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func parseChamber(r *http.Request) (chamber.Chamber, error) {
	var chamber chamber.Chamber
	err := json.NewDecoder(r.Body).Decode(&chamber)

	return chamber, errors.Wrap(err, "could not decode chamber from request body")
}
