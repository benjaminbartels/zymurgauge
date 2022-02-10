package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/settings"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type SettingsHandler struct {
	SettingsRepo settings.Repo
	UpdateChan   chan settings.Settings
}

func (h *SettingsHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	c, err := h.SettingsRepo.Get()
	if err != nil {
		return errors.Wrap(err, "could not get settings from controller")
	}

	if c == nil {
		return web.NewRequestError("settings not found", http.StatusNotFound)
	}

	if err := web.Respond(ctx, w, c, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *SettingsHandler) Save(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	s, err := parseSettings(r)
	if err != nil {
		return errors.Wrap(err, "could not parse settings")
	}

	if err := h.SettingsRepo.Save(&s); err != nil {
		return errors.Wrap(err, "could not save settings to controller")
	}

	if err := web.Respond(ctx, w, s, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	if h.UpdateChan != nil {
		h.UpdateChan <- s
	}

	return nil
}

func parseSettings(r *http.Request) (settings.Settings, error) {
	var settings settings.Settings
	err := json.NewDecoder(r.Body).Decode(&settings)

	return settings, errors.Wrap(err, "could not decode settings from request body")
}
