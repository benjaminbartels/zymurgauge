package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// Handle handles the incoming http request
func (h *ChamberHandler) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case web.GET:
		return h.get(ctx, w, r)
	case web.POST:
		return h.post(w, r)
	case web.DELETE:
		return h.delete(r)
	default:
		return web.ErrMethodNotAllowed
	}
}

func (h *ChamberHandler) get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var head string
	head, r.URL.Path = web.ShiftPath(r.URL.Path)
	if head == "" {
		return h.getAll(ctx, w)
	}
	mac, err := url.QueryUnescape(head)
	if err != nil {
		return web.ErrBadRequest
	}
	head, r.URL.Path = web.ShiftPath(r.URL.Path)
	if head == "events" {
		return h.getEvents(w, mac)
	} else if head == "" {
		return h.getOne(ctx, w, mac)
	} else {
		return web.ErrBadRequest
	}

}

func (h *ChamberHandler) getOne(ctx context.Context, w http.ResponseWriter, mac string) error {
	if chamber, err := h.repo.Get(mac); err != nil {
		return err
	} else if chamber == nil {
		return web.ErrNotFound
	} else {
		return web.Respond(ctx, w, chamber, http.StatusOK)
	}
}

func (h *ChamberHandler) getAll(ctx context.Context, w http.ResponseWriter) error {
	chambers, err := h.repo.GetAll()
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, chambers, http.StatusOK)
}

func (h *ChamberHandler) getEvents(w http.ResponseWriter, mac string) error {
	f, ok := w.(http.Flusher) //ToDo: comment this mess
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
		if _, err := fmt.Fprint(w, msg); err != nil {
			return err // ToDo: Test this
		}
		h.logger.Printf("Sending: %s\n", msg)
		f.Flush()
	}
	return nil
}

func (h *ChamberHandler) post(w io.Writer, r *http.Request) error {
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
	_, err = w.Write(b)
	return err
}

func (h *ChamberHandler) delete(r *http.Request) error {
	if r.URL.Path == "" {
		return web.ErrBadRequest
	}

	mac, err := url.QueryUnescape(r.URL.Path)
	if err != nil {
		return err
	}

	return h.repo.Delete(mac)
}

func parseChamber(r *http.Request) (internal.Chamber, error) {
	var chamber internal.Chamber
	err := json.NewDecoder(r.Body).Decode(&chamber)
	return chamber, err
}
