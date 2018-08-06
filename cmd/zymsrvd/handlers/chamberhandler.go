package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/pubsub"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

// ChamberHandler is the http handler for API calls to manage Chambers
type ChamberHandler struct {
	repo   *database.ChamberRepo
	pubSub *pubsub.PubSub
	logger log.Logger
}

// NewChamberHandler instantiates a ChamberHandler
func NewChamberHandler(repo *database.ChamberRepo, pubSub *pubsub.PubSub, logger log.Logger) *ChamberHandler {
	return &ChamberHandler{
		repo:   repo,
		pubSub: pubSub,
		logger: logger,
	}
}

// GetAll handles a GET request for all Chambers
func (h *ChamberHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p map[string]string) error {
	if chambers, err := h.repo.GetAll(); err != nil {
		return err
	} else {
		web.Respond(ctx, w, chambers, http.StatusOK)
	}
	return nil
}

// GetOne handles a GET request for a specific Chamber whose mac address matched the provided mac
func (h *ChamberHandler) GetOne(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p map[string]string) error {
	mac := p["mac"]
	if chamber, err := h.repo.Get(mac); err != nil {
		return err
	} else if chamber == nil {
		return web.ErrNotFound
	} else {
		web.Respond(ctx, w, chamber, http.StatusOK)
	}
	return nil
}

// GetEvents handles a GET request to listen for web events for a specific Chamber
// whose mac address matched the provided mac
func (h *ChamberHandler) GetEvents(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p map[string]string) error {
	mac := p["mac"]

	f, ok := w.(http.Flusher)
	if !ok {
		return web.ErrInternal
	}
	ch := h.pubSub.Subscribe(mac)
	h.logger.Printf("Added client [%s]\n", mac)
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		h.pubSub.Unsubscribe(ch)
		h.logger.Printf("Removed client %s channel\n", mac)
	}()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	for {
		c, ok := <-ch
		if !ok {
			break
		}
		msg := fmt.Sprintf("data: %s\n", c)
		h.logger.Printf("Sending: %s\n", msg)
		f.Flush()
	}

	return nil
}

// Post handles the POST request to create or update a Chamber
func (h *ChamberHandler) Post(ctx context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) error {
	chamber, err := parseChamber(r)
	if err != nil {
		return err
	}

	if err := h.repo.Save(&chamber); err != nil {
		return err
	} else {
		web.Respond(ctx, w, chamber, http.StatusOK)
	}
	return nil
}

// Delete handles the DELETE request to delete a Chamber
func (h *ChamberHandler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) error {
	mac := p["mac"]

	if err := h.repo.Delete(mac); err != nil {
		return err
	}

	web.Respond(ctx, w, nil, http.StatusOK)
	return nil
}

// parseChamber decodes the specified Chamber into JSON
func parseChamber(r *http.Request) (internal.Chamber, error) {
	var chamber internal.Chamber
	err := json.NewDecoder(r.Body).Decode(&chamber)
	return chamber, err
}
