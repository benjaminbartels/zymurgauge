package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/device/raspberrypi"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"periph.io/x/periph/experimental/host/netlink"
)

const (
	masterID = 0o01
)

type ThermometersHandler struct{}

func (h *ThermometersHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params) error {
	oneBus, err := netlink.New(masterID)
	if err != nil {
		log.Printf("Could not open Netlink host: %v", err)
	}

	ids, err := raspberrypi.GetThermometerIDs(oneBus)
	if err != nil {
		return errors.Wrap(err, "could not get all chambers from repository")
	}

	if err = web.Respond(ctx, w, ids, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}
