package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

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

func (h *ChamberHandler) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

func (h *ChamberHandler) handleGet(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var head string
	head, r.URL.Path = web.ShiftPath(r.URL.Path)
	if head == "" {
		return h.handleGetAll(ctx, w)
	} else {
		mac, err := url.QueryUnescape(head)
		if err != nil {
			return web.ErrBadRequest
		}
		head, r.URL.Path = web.ShiftPath(r.URL.Path)
		if head == "events" {
			return h.handleGetEvents(ctx, w, mac)
		} else if head == "" {
			return h.handleGetOne(ctx, w, mac)
		} else {
			return web.ErrBadRequest
		}
	}
}

func (h *ChamberHandler) handleGetOne(ctx context.Context, w http.ResponseWriter, mac string) error {
	if chamber, err := h.repo.Get(mac); err != nil {
		return err
	} else if chamber == nil {
		return web.ErrNotFound
	} else {
		return web.Respond(ctx, w, chamber, http.StatusOK)
	}
}

func (h *ChamberHandler) handleGetAll(ctx context.Context, w http.ResponseWriter) error {
	if chambers, err := h.repo.GetAll(); err != nil {
		return err
	} else {
		return web.Respond(ctx, w, chambers, http.StatusOK)
	}
}

func (h *ChamberHandler) handleGetEvents(ctx context.Context, w http.ResponseWriter, mac string) error {
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
		fmt.Fprint(w, msg)
		h.logger.Printf("Sending: %s\n", msg)
		f.Flush()
	}
	return nil
}

func (h *ChamberHandler) handlePost(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	chamber, err := parseChamber(r)
	if err != nil {
		return err
	}
	if err = h.repo.Save(&chamber); err != nil {
		return err
	}
	b, err := json.Marshal(chamber)
	if err != nil {
		return err
	}
	h.pubSub.Send(chamber.MacAddress, b)
	if _, err = w.Write(b); err != nil {
		return err
	}
	return nil
}

func (h *ChamberHandler) handleDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path == "" {
		return web.ErrBadRequest
	}

	if mac, err := url.QueryUnescape(r.URL.Path); err != nil {
		return err
	} else {
		if err := h.repo.Delete(mac); err != nil {
			return err
		}
	}
	return nil
}

func parseChamber(r *http.Request) (internal.Chamber, error) {
	var chamber internal.Chamber
	err := json.NewDecoder(r.Body).Decode(&chamber)
	return chamber, err
}
