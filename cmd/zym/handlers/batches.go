package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/batch"
	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type BatchesHandler struct {
	Service brewfather.Service
}

func (h *BatchesHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params,
) error {
	batchSummaries, err := h.Service.GetAllBatchSummaries(ctx)
	if err != nil {
		return errors.Wrap(err, "could not get all batches from repository")
	}

	summaries := batch.ConvertSummaries(batchSummaries)

	if err = web.Respond(ctx, w, summaries, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *BatchesHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	batchDetail, err := h.Service.GetBatchDetail(ctx, id)
	if err != nil {
		if errors.Is(err, brewfather.ErrNotFound) {
			return web.NewRequestError(fmt.Sprintf("batch '%s' not found", id), http.StatusNotFound)
		}

		return errors.Wrap(err, "could not get batch from repository")
	}

	if batchDetail == nil {
		return web.NewRequestError(fmt.Sprintf("batch '%s' not found", id), http.StatusNotFound)
	}

	detail := batch.ConvertDetail(batchDetail)

	if err = web.Respond(ctx, w, detail, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}
