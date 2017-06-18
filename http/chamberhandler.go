package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/benjaminbartels/zymurgauge"
	"github.com/benjaminbartels/zymurgauge/gpio"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// ChamberHandler is the HTTP resource handler for Chambers.
type ChamberHandler struct {
	*httprouter.Router
	ChamberService     zymurgauge.ChamberService
	logger             *logrus.Logger
	newSubscribers     chan subscriber
	closingSubscribers chan subscriber
}

// NewChamberHandler instantiates a new ChamberHandler
func NewChamberHandler(l *logrus.Logger) *ChamberHandler {
	h := &ChamberHandler{
		Router:             httprouter.New(),
		logger:             l,
		newSubscribers:     make(chan subscriber),
		closingSubscribers: make(chan subscriber),
	}
	h.GET("/api/chambers/:mac", h.handleGet)
	h.GET("/api/chambers/:mac/*option", h.handleGetEvents)
	h.POST("/api/chambers", h.handlePost)

	return h
}

// handleGet handles the HTTP GET request for Chambers by MAC address
func (h *ChamberHandler) handleGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	mac, err := url.QueryUnescape(ps.ByName("mac"))
	if err != nil {
		handleError(w, err, http.StatusBadRequest, h.logger) // ToDo: is this valid
	}

	f, err := h.ChamberService.Get(mac)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError, h.logger)
	} else if f == nil {
		notFound(w, h.logger)
	} else {
		encode(w, &getChamberResponse{Chamber: f}, h.logger)
	}
}

// handleGetEvents handles the HTTP GET request for Server Side Events
func (h *ChamberHandler) handleGetEvents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if ps.ByName("option") != "/events" {
		notFound(w, h.logger)
	}

	// Ensure support for flushing.
	f, ok := w.(http.Flusher)
	if !ok {
		handleError(w, zymurgauge.ErrInternal, http.StatusInternalServerError, h.logger) // ToDo: Make error better
		return
	}

	mac, err := url.QueryUnescape(ps.ByName("mac"))
	if err != nil {
		handleError(w, err, http.StatusBadRequest, h.logger) // ToDo: is this valid
	}

	// Create a new subscriber
	sub := subscriber{
		mac: mac,
		ch:  make(chan zymurgauge.Chamber),
	}

	h.newSubscribers <- sub

	// Listen for the the closing of the http connection
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify

		// Remove subscriber
		h.closingSubscribers <- sub
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {

		// Read from our messageChan
		c, ok := <-sub.ch

		if !ok {
			// Channel is closed, client disconnected
			break
		}

		data, err := json.Marshal(c)
		if err != nil {
			handleError(w, zymurgauge.ErrInvalidJSON, http.StatusInternalServerError, h.logger)
			return
		}

		msg := fmt.Sprintf("data: %s\n", string(data))
		fmt.Fprint(w, msg)

		h.logger.Debugf("Sending: %s", msg)

		if err != nil {
			handleError(w, zymurgauge.ErrInvalidJSON, http.StatusInternalServerError, h.logger)
			return
		}

		f.Flush()
	}
}

// handlePost handles the HTTP POST request for the Chamber resource
func (h *ChamberHandler) handlePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req postChamberRequest
	req.Chamber = &zymurgauge.Chamber{
		Controller: &gpio.Thermostat{},
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, zymurgauge.ErrInvalidJSON, http.StatusBadRequest, h.logger)
		return
	}
	f := req.Chamber
	f.ModTime = time.Time{}

	switch err := h.ChamberService.Save(f); err {
	case nil:
		encode(w, &postChamberResponse{Chamber: f}, h.logger)
	case zymurgauge.ErrChamberRequired:
		handleError(w, err, http.StatusBadRequest, h.logger)
	default:
		handleError(w, err, http.StatusInternalServerError, h.logger)
	}
}

// Start instructs the ChamberHandler to start processing GET requests for events
func (h *ChamberHandler) Start() { //ToDo: Rename?

	for {

		select {

		case s := <-h.newSubscribers:

			err := h.ChamberService.Subscribe(s.mac, s.ch)
			if err != nil {
				h.logger.Error(err)
			}
			h.logger.Debugf("Added client %s", s.mac)

		case s := <-h.closingSubscribers:
			h.ChamberService.Unsubscribe(s.mac)
			h.logger.Debugf("Removed client %s", s.mac)
		}
	}
}

// getChamberResponse represents the http response for a get
type getChamberResponse struct {
	Chamber *zymurgauge.Chamber `json:"chamber,omitempty"`
	Err     string              `json:"err,omitempty"`
}

// postChamberRequest represents the http request for a post
type postChamberRequest struct {
	Chamber *zymurgauge.Chamber `json:"chamber,omitempty"`
}

// postChamberResponse represents the http response for a post
type postChamberResponse struct {
	Chamber *zymurgauge.Chamber `json:"chamber,omitempty"`
	Err     string              `json:"err,omitempty"`
}

// subscriber represents a client subscribe for Chamber update events
type subscriber struct {
	mac string
	ch  chan zymurgauge.Chamber
}
