package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type BatchesHandler struct {
	Service brewfather.Service
}

func (h *BatchesHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params) error {
	batches, err := h.Service.GetAll(ctx)
	if err != nil {
		return errors.Wrap(err, "could not get all batches from repository")
	}

	if err = web.Respond(ctx, w, batches, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *BatchesHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	batch, err := h.Service.Get(ctx, id)
	if err != nil {
		if errors.Is(err, brewfather.ErrNotFound) {
			return web.NewRequestError(fmt.Sprintf("batch '%s' not found", id), http.StatusNotFound)
		}

		return errors.Wrap(err, "could not get batch from repository")
	}

	if batch == nil {
		return web.NewRequestError(fmt.Sprintf("batch '%s' not found", id), http.StatusNotFound)
	}

	if err = web.Respond(ctx, w, batch, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}
