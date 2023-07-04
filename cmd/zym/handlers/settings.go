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

func (h *SettingsHandler) Get(ctx context.Context, w http.ResponseWriter, _ *http.Request, _ httprouter.Params) error {
	s, err := h.SettingsRepo.Get()
	if err != nil {
		return errors.Wrap(err, "could not get settings from repository")
	}

	if s == nil {
		return web.NewRequestError("settings not found", http.StatusNotFound)
	}

	if err := web.Respond(ctx, w, s.AppSettings, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *SettingsHandler) Save(ctx context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	appSettings, err := parseAppSettings(r)
	if err != nil {
		return errors.Wrap(err, "could not parse settings")
	}

	s, err := h.SettingsRepo.Get()
	if err != nil {
		return errors.Wrap(err, "could not get settings from repository")
	}

	s.AppSettings = appSettings

	if err := h.SettingsRepo.Save(s); err != nil {
		return errors.Wrap(err, "could not save settings to repository")
	}

	if err := web.Respond(ctx, w, s, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	if h.UpdateChan != nil {
		h.UpdateChan <- *s
	}

	return nil
}

func parseAppSettings(r *http.Request) (settings.AppSettings, error) {
	var settings settings.AppSettings
	err := json.NewDecoder(r.Body).Decode(&settings)

	return settings, errors.Wrap(err, "could not decode settings from request body")
}
