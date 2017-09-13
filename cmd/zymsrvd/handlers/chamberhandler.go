package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/platform/pubsub"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

type ChamberHandler struct {
	repo   *database.ChamberRepo
	pubSub *pubsub.PubSub
}

func NewChamberHandler(repo *database.ChamberRepo, pubSub *pubsub.PubSub) *ChamberHandler {

	return &ChamberHandler{
		repo:   repo,
		pubSub: pubSub,
	}
}

func (h *ChamberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGet(w, r)
	case "POST":
		h.handlePost(w, r)
	default:
		web.HandleError(w, web.ErrNotFound)
	}
}

func (h *ChamberHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	if head == "" {
		h.handleGetAll(w)
	} else {

		mac, err := url.QueryUnescape(head)
		if err != nil {
			web.HandleError(w, web.ErrBadRequest)
		}

		head, r.URL.Path = shiftPath(r.URL.Path)

		if head == "events" {
			h.handleGetEvents(w, mac)
		} else {
			h.handleGetOne(w, mac)
		}
	}

}

func (h *ChamberHandler) handleGetOne(w http.ResponseWriter, mac string) {
	if chamber, err := h.repo.Get(mac); err != nil {
		web.HandleError(w, err)
	} else if chamber == nil {
		web.HandleError(w, web.ErrNotFound)
	} else {
		web.Encode(w, &chamber)
	}

}

func (h *ChamberHandler) handleGetAll(w http.ResponseWriter) {
	if chambers, err := h.repo.GetAll(); err != nil {
		web.HandleError(w, err)
	} else {
		web.Encode(w, chambers)
	}
}

func (h *ChamberHandler) handleGetEvents(w http.ResponseWriter, mac string) {

	f, ok := w.(http.Flusher)
	if !ok {
		web.HandleError(w, web.ErrInternal)
		return
	}

	ch := h.pubSub.Subscribe(mac)

	fmt.Printf("Added client [%s]\n", mac)

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify

		h.pubSub.Unsubscribe(ch)
		fmt.Printf("Removed client %s channel\n", mac)
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

		fmt.Printf("Sending: %s\n", msg)

		f.Flush()
	}
}

func (h *ChamberHandler) handlePost(w http.ResponseWriter, r *http.Request) {

	chamber, err := parseChamber(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	if err = h.repo.Save(&chamber); err != nil {
		web.HandleError(w, err)
		return
	}

	b, err := json.Marshal(chamber)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	h.pubSub.Send(chamber.MacAddress, b)

	if _, err = w.Write(b); err != nil {
		web.HandleError(w, err)
		return
	}

}

func parseChamber(r *http.Request) (internal.Chamber, error) {
	var chamber internal.Chamber
	err := json.NewDecoder(r.Body).Decode(&chamber)
	return chamber, err
}
