package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/platform/app"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/pubsub"
)

// ChamberHandler is the http handler for API calls to manage Chambers
type ChamberHandler struct {
	app.Handler
	repo   *database.ChamberRepo
	pubSub *pubsub.PubSub
	logger log.Logger
}

// NewChamberHandler instantiates a ChamberHandler
func NewChamberHandler(repo *database.ChamberRepo, pubSub *pubsub.PubSub, logger log.Logger) *ChamberHandler {
	return &ChamberHandler{
		Handler: app.Handler{Logger: logger},
		repo:    repo,
		pubSub:  pubSub,
		logger:  logger,
	}
}

// ServeHTTP calls f(w, r).
func (h *ChamberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case app.GET:
		h.handleGet(w, r)
	case app.POST:
		h.handlePost(w, r)
	case app.DELETE:
		h.handleDelete(w, r)
	default:
		h.HandleError(w, app.ErrNotFound)
	}
}

func (h *ChamberHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = h.ShiftPath(r.URL.Path)
	if head == "" {
		h.handleGetAll(w)
	} else {
		mac, err := url.QueryUnescape(head)
		if err != nil {
			h.HandleError(w, app.ErrBadRequest)
		}
		head, r.URL.Path = h.ShiftPath(r.URL.Path)
		if head == "events" {
			h.handleGetEvents(w, mac)
		} else if head == "" {
			h.handleGetOne(w, mac)
		} else {
			h.HandleError(w, app.ErrBadRequest)
		}
	}
}

func (h *ChamberHandler) handleGetOne(w http.ResponseWriter, mac string) {
	if chamber, err := h.repo.Get(mac); err != nil {
		h.HandleError(w, err)
	} else if chamber == nil {
		h.HandleError(w, app.ErrNotFound)
	} else {
		h.Encode(w, &chamber)
	}
}

func (h *ChamberHandler) handleGetAll(w http.ResponseWriter) {
	if chambers, err := h.repo.GetAll(); err != nil {
		h.HandleError(w, err)
	} else {
		h.Encode(w, chambers)
	}
}

func (h *ChamberHandler) handleGetEvents(w http.ResponseWriter, mac string) {
	f, ok := w.(http.Flusher)
	if !ok {
		h.HandleError(w, app.ErrInternal)
		return
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
}

func (h *ChamberHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	chamber, err := parseChamber(r)
	if err != nil {
		h.HandleError(w, err)
		return
	}
	if err = h.repo.Save(&chamber); err != nil {
		h.HandleError(w, err)
		return
	}
	b, err := json.Marshal(chamber)
	if err != nil {
		h.HandleError(w, err)
		return
	}
	h.pubSub.Send(chamber.MacAddress, b)
	if _, err = w.Write(b); err != nil {
		h.HandleError(w, err)
		return
	}
}

func (h *ChamberHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "" {
		if id, err := strconv.ParseUint(r.URL.Path, 10, 64); err != nil {
			h.HandleError(w, app.ErrBadRequest)
		} else {
			if err := h.repo.Delete(id); err != nil {
				h.HandleError(w, err)
			}
		}
		return
	}
	h.HandleError(w, app.ErrBadRequest)
}

func parseChamber(r *http.Request) (internal.Chamber, error) {
	var chamber internal.Chamber
	err := json.NewDecoder(r.Body).Decode(&chamber)
	return chamber, err
}
