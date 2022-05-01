package handlers

import (
	"context"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/device/onewire"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type ThermometersHandler struct {
	DevicePath string
}

func (h *ThermometersHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params,
) error {
	ids, err := onewire.GetIDs(h.DevicePath, onewire.Ds18b20Prefix)
	if err != nil {
		return errors.Wrap(err, "could not get all thermometers ids from onewire bus")
	}

	if err = web.Respond(ctx, w, ids, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}
