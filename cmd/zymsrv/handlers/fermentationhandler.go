package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
)

// FermentationHandler is the http handler for API calls to manage Fermentations.
type FermentationHandler struct {
	fermRepo    *storage.FermentationRepo
	changeRepo  *storage.TemperatureChangeRepo
	chamberRepo *storage.ChamberRepo
}

// NewFermentationHandler instantiates a FermentationHandler.
func NewFermentationHandler(fermRepo *storage.FermentationRepo, changeRepo *storage.TemperatureChangeRepo,
	chamberRepo *storage.ChamberRepo) *FermentationHandler {
	return &FermentationHandler{
		fermRepo:    fermRepo,
		changeRepo:  changeRepo,
		chamberRepo: chamberRepo,
	}
}

// Handle handles the incoming http request.
func (h *FermentationHandler) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case web.GET:
		return h.get(ctx, w, r)
	case web.POST:
		return h.post(ctx, w, r)
	case web.DELETE:
		return h.delete(ctx, w, r)
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

	switch {
	case head == "temperaturechanges":
		start := time.Now().AddDate(0, 0, -1).UTC()
		//nolint:gomnd
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
	case head == "":
		return h.getOne(ctx, w, id)
	default:
		return web.ErrBadRequest
	}
}

func (h *FermentationHandler) getOne(ctx context.Context, w http.ResponseWriter, id uint64) error {
	fermentation, err := h.fermRepo.Get(id)

	switch {
	case err != nil:
		return err
	case fermentation == nil:
		return web.ErrNotFound
	default:
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

	id, err := strconv.ParseUint(head, 10, 64)
	if err != nil {
		return web.ErrBadRequest
	}

	head, r.URL.Path = web.ShiftPath(r.URL.Path)

	switch head {
	case "temperaturechanges":
		return h.postTemperatureChange(ctx, w, r)
	case "start":
		return h.postStart(ctx, w, id)
	case "stop":
		return h.postStop(ctx, w, id)
	}

	return web.ErrBadRequest
}

func (h *FermentationHandler) postFermentation(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

func (h *FermentationHandler) postTemperatureChange(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

func (h *FermentationHandler) postStart(ctx context.Context, w http.ResponseWriter, id uint64) error {
	fermentation, err := h.fermRepo.Get(id)
	if err != nil {
		return err
	}

	chamber, err := h.chamberRepo.Get(fermentation.Chamber.MacAddress)
	if err != nil {
		return err
	}

	chamber.CurrentFermentationID = fermentation.ID

	if err := h.chamberRepo.Save(chamber); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}

func (h *FermentationHandler) postStop(ctx context.Context, w http.ResponseWriter, id uint64) error {
	fermentation, err := h.fermRepo.Get(id)
	if err != nil {
		return err
	}

	chamber, err := h.chamberRepo.Get(fermentation.Chamber.MacAddress)
	if err != nil {
		return err
	}

	chamber.CurrentFermentationID = 0

	if err := h.chamberRepo.Save(chamber); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}

func (h *FermentationHandler) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var head string
	head, r.URL.Path = web.ShiftPath(r.URL.Path)

	if head == "" {
		return web.ErrBadRequest
	}

	id, err := strconv.ParseUint(head, 10, 64)
	if err != nil {
		return web.ErrBadRequest
	}

	if err := h.fermRepo.Delete(id); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}

func parseFermentation(r *http.Request) (storage.Fermentation, error) {
	var fermentation storage.Fermentation
	err := json.NewDecoder(r.Body).Decode(&fermentation)

	return fermentation, err
}

func parseTemperatureChange(r *http.Request) (internal.TemperatureChange, error) {
	var change internal.TemperatureChange
	err := json.NewDecoder(r.Body).Decode(&change)

	return change, err
}
