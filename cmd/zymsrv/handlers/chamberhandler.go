package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/pubsub"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
)

// ChamberHandler is the http handler for API calls to manage Chambers.
type ChamberHandler struct {
	repo   *storage.ChamberRepo
	pubSub *pubsub.PubSub
	logger log.Logger
}

// NewChamberHandler instantiates a ChamberHandler.
func NewChamberHandler(repo *storage.ChamberRepo, pubSub *pubsub.PubSub, logger log.Logger) *ChamberHandler {
	return &ChamberHandler{
		repo:   repo,
		pubSub: pubSub,
		logger: logger,
	}
}

// Handle handles the incoming http request.
func (h *ChamberHandler) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case web.GET:
		return h.get(ctx, w, r)
	case web.POST:
		return h.post(w, r)
	case web.DELETE:
		return h.delete(ctx, w, r)
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

	switch head {
	case "":
		return h.getOne(ctx, w, mac)
	default:
		return web.ErrBadRequest
	}
}

func (h *ChamberHandler) getOne(ctx context.Context, w http.ResponseWriter, mac string) error {
	chamber, err := h.repo.Get(mac)

	switch {
	case err != nil:
		return err
	case chamber == nil:
		return web.ErrNotFound
	default:
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

func (h *ChamberHandler) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var head string
	head, r.URL.Path = web.ShiftPath(r.URL.Path)

	if head == "" {
		return web.ErrBadRequest
	}

	if err := h.repo.Delete(head); err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}

func parseChamber(r *http.Request) (storage.Chamber, error) {
	var chamber storage.Chamber
	err := json.NewDecoder(r.Body).Decode(&chamber)

	return chamber, err
}
