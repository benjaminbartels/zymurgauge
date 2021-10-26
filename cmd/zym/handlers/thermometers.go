package handlers

import (
	"context"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/device/raspberrypi"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type ThermometersHandler struct{}

func (h *ThermometersHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params) error {
	ids, err := raspberrypi.GetThermometerIDs()
	if err != nil {
		return errors.Wrap(err, "could not get all chambers from repository")
	}

	if err = web.Respond(ctx, w, ids, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}
