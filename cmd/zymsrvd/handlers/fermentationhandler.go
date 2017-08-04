package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

type FermentationHandler struct {
	repo *database.FermentationRepo
}

func NewFermentationHandler(repo *database.FermentationRepo) *FermentationHandler {
	return &FermentationHandler{
		repo: repo,
	}
}

func (h *FermentationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGet(w, r)
	case "POST":
		h.handlePost(w, r)
	default:
		web.HandleError(w, web.ErrNotFound)
	}
}

func (h *FermentationHandler) handleGet(w http.ResponseWriter, r *http.Request) {

	id, err := parseID(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	if fermentation, err := h.repo.Get(id); err != nil {
		web.HandleError(w, err)
	} else if fermentation == nil {
		web.HandleError(w, web.ErrNotFound)
	} else {
		web.Encode(w, &fermentation)
	}

}

func (h *FermentationHandler) handlePost(w http.ResponseWriter, r *http.Request) {

	fermentation, err := parseFermentation(r)
	if err != nil {
		web.HandleError(w, err)
		return
	}

	fmt.Println(fermentation)

	if err := h.repo.Save(&fermentation); err != nil {
		web.HandleError(w, err)
	} else {
		web.Encode(w, &fermentation)
	}
}

func parseFermentation(r *http.Request) (internal.Fermentation, error) {
	var fermentation internal.Fermentation
	err := json.NewDecoder(r.Body).Decode(&fermentation)
	return fermentation, err
}
