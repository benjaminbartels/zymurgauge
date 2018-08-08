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

// Handle handles the incoming http request
func (h *FermentationHandler) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case web.GET:
		return h.get(ctx, w, r)
	case web.POST:
		return h.post(ctx, w, r)
	case web.DELETE:
		return h.delete(r)
	default:
		return web.ErrMethodNotAllowed
	}
}

func (h *FermentationHandler) get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var head string
	head, r.URL.Path = web.ShiftPath(r.URL.Path)
	if head == "" {
		return h.getAll(ctx, w)
	}
	id, err := strconv.ParseUint(head, 10, 64)
	if err != nil {
		return web.ErrBadRequest
	}
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
		return h.getTemperatureChanges(ctx, w, id, start, end)
	} else if head == "" {
		return h.getOne(ctx, w, id)
	} else {
		return web.ErrBadRequest
	}
}

func (h *FermentationHandler) getOne(ctx context.Context, w http.ResponseWriter, id uint64) error {
	if fermentation, err := h.fermRepo.Get(id); err != nil {
		return err
	} else if fermentation == nil {
		return web.ErrNotFound
	} else {
		return web.Respond(ctx, w, fermentation, http.StatusOK)
	}
}

func (h *FermentationHandler) getAll(ctx context.Context, w http.ResponseWriter) error {
	fermentations, err := h.fermRepo.GetAll()
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, fermentations, http.StatusOK)
}

func (h *FermentationHandler) getTemperatureChanges(ctx context.Context, w http.ResponseWriter, id uint64,
	start, end time.Time) error {
	changes, err := h.changeRepo.GetRangeByFermentationID(id, start, end)
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, changes, http.StatusOK)
}

func (h *FermentationHandler) post(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var head string
	head, r.URL.Path = web.ShiftPath(r.URL.Path)
	if head == "" {
		return h.postFermentation(ctx, w, r)
	}
	_, err := strconv.ParseUint(head, 10, 64)
	if err != nil {
		return web.ErrBadRequest
	}
	head, r.URL.Path = web.ShiftPath(r.URL.Path)
	if head == "temperaturechanges" {
		return h.postTemperatureChange(ctx, w, r)
	}
	return web.ErrBadRequest
}

func (h *FermentationHandler) postFermentation(ctx context.Context, w http.ResponseWriter,
	r *http.Request) error {
	fermentation, err := parseFermentation(r)
	if err != nil {
		return err
	}
	err = h.fermRepo.Save(&fermentation)
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, fermentation, http.StatusOK)
}

func (h *FermentationHandler) postTemperatureChange(ctx context.Context, w http.ResponseWriter,
	r *http.Request) error {
	change, err := parseTemperatureChange(r)
	if err != nil {
		return err
	}
	err = h.changeRepo.Save(&change)
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, change, http.StatusOK)
}

func (h *FermentationHandler) delete(r *http.Request) error {
	if r.URL.Path == "" {
		return web.ErrBadRequest
	}

	id, err := strconv.ParseUint(r.URL.Path, 10, 64)
	if err != nil {
		return web.ErrBadRequest
	}
	if err := h.fermRepo.Delete(id); err != nil {
		return err
	}
	// ToDo: delete temperaturechanges
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
