package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

// FermentationHandler is the http handler for API calls to manage Fermentations
type FermentationHandler struct {
	fermRepo   *database.FermentationRepo
	changeRepo *database.TemperatureChangeRepo
}

// NewFermentationHandler instantiates a FermentationHandler
func NewFermentationHandler(fermRepo *database.FermentationRepo,
	changeRepo *database.TemperatureChangeRepo) *FermentationHandler {
	return &FermentationHandler{
		fermRepo:   fermRepo,
		changeRepo: changeRepo,
	}
}

func (h *FermentationHandler) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case web.GET:
		return h.handleGet(ctx, w, r)
	case web.POST:
		return h.handlePost(ctx, w, r)
	case web.DELETE:
		return h.handleDelete(ctx, w, r)
	default:
		return web.ErrMethodNotAllowed
	}
}

func (h *FermentationHandler) handleGet(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var head string
	head, r.URL.Path = web.ShiftPath(r.URL.Path)
	if head == "" {
		return h.handleGetAll(ctx, w)
	} else {
		if id, err := strconv.ParseUint(head, 10, 64); err != nil {
			return web.ErrBadRequest
		} else {
			head, r.URL.Path = web.ShiftPath(r.URL.Path)
			if head == "temperaturechanges" {
				start := time.Time{}
				end := time.Unix(1<<63-62135596801, 999999999).UTC()
				if startParam, ok := r.URL.Query()["start"]; ok {
					start, err = time.Parse(time.RFC3339, startParam[0])
					if err != nil {
						return err
					}
				}
				if endParam, ok := r.URL.Query()["end"]; ok {
					end, err = time.Parse(time.RFC3339, endParam[0])
					if err != nil {
						return err
					}
				}
				return h.handleGetTemperatureChanges(ctx, w, id, start, end)
			} else if head == "" {
				return h.handleGetOne(ctx, w, id)
			} else {
				return web.ErrBadRequest
			}
		}
	}
}

func (h *FermentationHandler) handleGetOne(ctx context.Context, w http.ResponseWriter, id uint64) error {
	if fermentation, err := h.fermRepo.Get(id); err != nil {
		return err
	} else if fermentation == nil {
		return web.ErrNotFound
	} else {
		return web.Respond(ctx, w, fermentation, http.StatusOK)
	}
}

func (h *FermentationHandler) handleGetAll(ctx context.Context, w http.ResponseWriter) error {
	if fermentations, err := h.fermRepo.GetAll(); err != nil {
		return err
	} else {
		return web.Respond(ctx, w, fermentations, http.StatusOK)
	}
}

func (h *FermentationHandler) handleGetTemperatureChanges(ctx context.Context, w http.ResponseWriter, id uint64,
	start, end time.Time) error {
	if changes, err := h.changeRepo.GetRangeByFermentationID(id, start, end); err != nil {
		return err
	} else {
		return web.Respond(ctx, w, changes, http.StatusOK)
	}
}

func (h *FermentationHandler) handlePost(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var head string
	head, r.URL.Path = web.ShiftPath(r.URL.Path)
	if head == "" {
		return h.handlePostFermentation(ctx, w, r)
	} else {
		if _, err := strconv.ParseUint(head, 10, 64); err != nil {
			return web.ErrBadRequest
		} else {
			head, r.URL.Path = web.ShiftPath(r.URL.Path)
			if head == "temperaturechanges" {
				return h.handlePostTemperatureChange(ctx, w, r)
			} else {
				return web.ErrBadRequest
			}
		}
	}
}

func (h *FermentationHandler) handlePostFermentation(ctx context.Context, w http.ResponseWriter,
	r *http.Request) error {
	fermentation, err := parseFermentation(r)
	if err != nil {
		return err
	}
	if err := h.fermRepo.Save(&fermentation); err != nil {
		return err
	} else {
		return web.Respond(ctx, w, fermentation, http.StatusOK)
	}
}

func (h *FermentationHandler) handlePostTemperatureChange(ctx context.Context, w http.ResponseWriter,
	r *http.Request) error {
	change, err := parseTemperatureChange(r)
	if err != nil {
		return err
	}
	if err := h.changeRepo.Save(&change); err != nil {
		return err
	} else {
		return web.Respond(ctx, w, change, http.StatusOK)
	}
}

func (h *FermentationHandler) handleDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "" {
		return web.ErrBadRequest
	}

	if id, err := strconv.ParseUint(r.URL.Path, 10, 64); err != nil {
		return web.ErrBadRequest
	} else {
		if err := h.fermRepo.Delete(id); err != nil {
			return err
		}

		// ToDo: delete temperaturechanges

	}
	return nil
}

func parseFermentation(r *http.Request) (internal.Fermentation, error) {
	var fermentation internal.Fermentation
	err := json.NewDecoder(r.Body).Decode(&fermentation)
	return fermentation, err
}

func parseTemperatureChange(r *http.Request) (internal.TemperatureChange, error) {
	var change internal.TemperatureChange
	err := json.NewDecoder(r.Body).Decode(&change)
	return change, err
}
