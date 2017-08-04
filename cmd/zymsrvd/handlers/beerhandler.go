package handlers

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

type BeerHandler struct {
	repo *database.BeerRepo
}

func NewBeerHandler(repo *database.BeerRepo) *BeerHandler {
	return &BeerHandler{
		repo: repo,
	}
}

func (h *BeerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGet(w, r)
	case "POST":
		h.handlePost(w, r)
	default:
		web.HandleError(w, web.ErrNotFound)
	}
}

func (h *BeerHandler) handleGet(w http.ResponseWriter, r *http.Request) {

	id, err := parseID(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	if beer, err := h.repo.Get(id); err != nil {
		web.HandleError(w, err)
	} else if beer == nil {
		web.HandleError(w, web.ErrNotFound)
	} else {
		web.Encode(w, &beer)
	}

}

func (h *BeerHandler) handlePost(w http.ResponseWriter, r *http.Request) {

	beer, err := parseBeer(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	fmt.Println(beer)

	if err := h.repo.Save(&beer); err != nil {
		web.HandleError(w, err)
	} else {
		web.Encode(w, &beer)
	}
}

func parseBeer(r *http.Request) (internal.Beer, error) {
	var beer internal.Beer
	err := json.NewDecoder(r.Body).Decode(&beer)
	return beer, err
}
