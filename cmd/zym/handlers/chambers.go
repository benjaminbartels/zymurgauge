package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ChambersHandler struct {
	ChamberController chamber.Controller
	Logger            *logrus.Logger
}

func (h *ChambersHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params,
) error {
	chambers, err := h.ChamberController.GetAll()

	for _, c := range chambers {
		c.RefreshReadings()
	}

	if err != nil {
		return errors.Wrap(err, "could not get all chambers from controller")
	}

	if err := web.Respond(ctx, w, chambers, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *ChambersHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	c, err := h.ChamberController.Get(id)
	if err != nil {
		return errors.Wrapf(err, "could not get chamber %s from controller", id)
	}

	if c == nil {
		return web.NewRequestError(fmt.Sprintf("chamber '%s' not found", id), http.StatusNotFound)
	}

	c.RefreshReadings()

	if err := web.Respond(ctx, w, c, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *ChambersHandler) Save(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	c, err := parseChamber(r)
	if err != nil {
		return errors.Wrap(err, "could not parse chamber")
	}

	if err := h.ChamberController.Save(&c); err != nil {
		var cfgError *chamber.InvalidConfigurationError

		switch {
		case errors.As(err, &cfgError):
			var b strings.Builder

			b.WriteString(fmt.Sprintf("%s: ", cfgError.Error()))

			for _, problem := range cfgError.Problems() {
				b.WriteString(fmt.Sprintf("%s, ", problem))
			}

			return web.NewRequestError(strings.TrimSuffix(b.String(), ", "), http.StatusBadRequest)

		case errors.Is(err, chamber.ErrFermenting):
			return web.NewRequestError("fermentation is in progress", http.StatusBadRequest)
		default:
			return errors.Wrap(err, "could not save chamber to controller")
		}
	}

	if err := web.Respond(ctx, w, c, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *ChambersHandler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params,
) error {
	id := p.ByName("id")

	if err := h.ChamberController.Delete(id); err != nil {
		switch {
		case errors.Is(err, chamber.ErrNotFound):
			return web.NewRequestError(fmt.Sprintf("chamber '%s' not found", id), http.StatusNotFound)
		case errors.Is(err, chamber.ErrFermenting):
			return web.NewRequestError("fermentation is in progress", http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "could not delete chamber %s from controller", id)
		}
	}

	if err := web.Respond(ctx, w, &Status{Message: "Success"}, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *ChambersHandler) Start(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params,
) error {
	id := p.ByName("id")

	step := r.URL.Query().Get("step")

	if err := h.ChamberController.StartFermentation(id, step); err != nil {
		switch {
		case errors.Is(err, chamber.ErrNotFound):
			return web.NewRequestError(fmt.Sprintf("chamber '%s' not found", id), http.StatusNotFound)
		case errors.Is(err, chamber.ErrInvalidStep):
			return web.NewRequestError(fmt.Sprintf("step '%s' is invalid for chamber '%s'", step, id), http.StatusBadRequest)
		case errors.Is(err, chamber.ErrNoCurrentBatch):
			return web.NewRequestError(fmt.Sprintf("chamber '%s' does not have a current batch", id), http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "could not start fermentation for chamber %s", id)
		}
	}

	if err := web.Respond(ctx, w, &Status{Message: "Success"}, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *ChambersHandler) Stop(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	if err := h.ChamberController.StopFermentation(id); err != nil {
		switch {
		case errors.Is(err, chamber.ErrNotFound):
			return web.NewRequestError(fmt.Sprintf("chamber '%s' not found", id), http.StatusNotFound)
		case errors.Is(err, chamber.ErrNotFermenting):
			return web.NewRequestError(fmt.Sprintf("chamber '%s' is not fermenting", id), http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "could not stop fermentation for chamber %s", id)
		}
	}

	if err := web.Respond(ctx, w, &Status{Message: "Success"}, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func parseChamber(r *http.Request) (chamber.Chamber, error) {
	var chamber chamber.Chamber
	err := json.NewDecoder(r.Body).Decode(&chamber)

	return chamber, errors.Wrap(err, "could not decode chamber from request body")
}
