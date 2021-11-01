package handlers

import (
	"context"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/device"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type ThermometersHandler struct {
	Repo device.ThermometerRepo
}

func (h *ThermometersHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params) error {
	ids, err := h.Repo.GetThermometerIDs()
	if err != nil {
		return errors.Wrap(err, "could not get all thermometers from repository")
	}

	if err = web.Respond(ctx, w, ids, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}
