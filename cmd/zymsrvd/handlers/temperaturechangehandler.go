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

// TemperatureChangeHandler is the http handler for API calls to manage TemperatureChange
type TemperatureChangeHandler struct {
	repo *database.TemperatureChangeRepo
}

// NewTemperatureChangeHandler instantiates a TemperatureChangeHandler
func NewTemperatureChangeHandler(repo *database.TemperatureChangeRepo) *TemperatureChangeHandler {
	return &TemperatureChangeHandler{
		repo: repo,
	}
}

// GetAll handles a GET request for all TemperatureChanges for a given FermentationID and date range
func (h *TemperatureChangeHandler) GetRange(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p map[string]string) error {

	id := p["id"]
	start := p["start"]
	end := p["end"]

	fermentationID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err //ToDo: error InvalidID
	}

	startTime := time.Time{}
	endTime := time.Unix(1<<63-62135596801, 999999999).UTC()

	startTime, err = time.Parse(time.RFC3339, start)
	if err != nil {
		return err //ToDo: error InvalidID

	}

	endTime, err = time.Parse(time.RFC3339, end)
	if err != nil {
		return err //ToDo: error InvalidID

	}

	if fermentations, err := h.repo.GetRangeByFermentationID(fermentationID, startTime, endTime); err != nil {
		return err
	} else {
		web.Respond(ctx, w, fermentations, http.StatusOK)
	}
	return nil
}

// Post handles the POST request to create or update a TemperatureChange
func (h *TemperatureChangeHandler) Post(ctx context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) error {
	change, err := parseTemperatureChange(r)
	if err != nil {
		return err
	}

	if err := h.repo.Save(&change); err != nil {
		return err
	}

	web.Respond(ctx, w, change, http.StatusOK)
	return nil
}

// parseTemperatureChange decodes the specified TemperatureChange into JSON
func parseTemperatureChange(r *http.Request) (internal.TemperatureChange, error) {
	var change internal.TemperatureChange
	err := json.NewDecoder(r.Body).Decode(&change)
	return change, err
}
