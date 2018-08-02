package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/benjaminbartels/zymurgauge/internal"
// 	"github.com/benjaminbartels/zymurgauge/internal/database"
// 	"github.com/benjaminbartels/zymurgauge/internal/platform/app"
// 	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
// )

// // FermentationHandler is the http handler for API calls to manage Fermentations
// type FermentationHandler struct {
// 	app.Handler
// 	fermRepo   *database.FermentationRepo
// 	changeRepo *database.TemperatureChangeRepo
// }

// // NewFermentationHandler instantiates a FermentationHandler
// func NewFermentationHandler(fermRepo *database.FermentationRepo, changeRepo *database.TemperatureChangeRepo,
// 	logger log.Logger) *FermentationHandler {
// 	return &FermentationHandler{
// 		Handler:    app.Handler{Logger: logger},
// 		fermRepo:   fermRepo,
// 		changeRepo: changeRepo,
// 	}
// }

// // ServeHTTP calls f(w, r).
// func (h *FermentationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case app.GET:
// 		h.handleGet(w, r)
// 	case app.POST:
// 		h.handlePost(w, r)
// 	case app.DELETE:
// 		h.handleDelete(w, r)
// 	default:
// 		h.HandleError(w, app.ErrNotFound)
// 	}
// }

// func (h *FermentationHandler) handleGet(w http.ResponseWriter, r *http.Request) {
// 	var head string
// 	head, r.URL.Path = h.ShiftPath(r.URL.Path)
// 	if head == "" {
// 		h.handleGetAll(w)
// 	} else {
// 		if id, err := strconv.ParseUint(head, 10, 64); err != nil {
// 			h.HandleError(w, app.ErrBadRequest)
// 		} else {
// 			head, r.URL.Path = h.ShiftPath(r.URL.Path)
// 			if head == "temperaturechanges" {
// 				start := time.Time{}
// 				end := time.Unix(1<<63-62135596801, 999999999).UTC()
// 				if startParam, ok := r.URL.Query()["start"]; ok {
// 					start, err = time.Parse(time.RFC3339, startParam[0])
// 					if err != nil {
// 						h.HandleError(w, err)
// 					}
// 				}
// 				if endParam, ok := r.URL.Query()["end"]; ok {
// 					end, err = time.Parse(time.RFC3339, endParam[0])
// 					if err != nil {
// 						h.HandleError(w, err)
// 					}
// 				}
// 				h.handleGetTemperatureChanges(w, id, start, end)
// 			} else if head == "" {
// 				h.handleGetOne(w, id)
// 			} else {
// 				h.HandleError(w, app.ErrBadRequest)
// 			}
// 		}
// 	}
// }

// func (h *FermentationHandler) handleGetOne(w http.ResponseWriter, id uint64) {
// 	if fermentation, err := h.fermRepo.Get(id); err != nil {
// 		h.HandleError(w, err)
// 	} else if fermentation == nil {
// 		h.HandleError(w, app.ErrNotFound)
// 	} else {
// 		h.Encode(w, &fermentation)
// 	}
// }

// func (h *FermentationHandler) handleGetAll(w http.ResponseWriter) {
// 	if fermentations, err := h.fermRepo.GetAll(); err != nil {
// 		h.HandleError(w, err)
// 	} else {
// 		fmt.Printf("%+v\n", fermentations)

// 		h.Encode(w, fermentations)
// 	}
// }

// func (h *FermentationHandler) handleGetTemperatureChanges(w http.ResponseWriter, id uint64, start, end time.Time) {
// 	if fermentations, err := h.changeRepo.GetRangeByFermentationID(id, start, end); err != nil {
// 		h.HandleError(w, err)
// 	} else {
// 		h.Encode(w, fermentations)
// 	}
// }

// func (h *FermentationHandler) handlePost(w http.ResponseWriter, r *http.Request) {
// 	var head string
// 	head, r.URL.Path = h.ShiftPath(r.URL.Path)
// 	if head == "" {

// 		h.handlePostFermentation(w, r)
// 	} else {
// 		if _, err := strconv.ParseUint(head, 10, 64); err != nil {
// 			h.HandleError(w, app.ErrBadRequest)
// 		} else {
// 			head, r.URL.Path = h.ShiftPath(r.URL.Path)
// 			if head == "temperaturechanges" {
// 				h.handlePostTemperatureChange(w, r)
// 			} else {
// 				h.HandleError(w, app.ErrBadRequest)
// 			}
// 		}
// 	}
// }

// func (h *FermentationHandler) handlePostFermentation(w http.ResponseWriter, r *http.Request) {
// 	fermentation, err := parseFermentation(r)
// 	if err != nil {
// 		h.HandleError(w, err)
// 		return
// 	}
// 	if err := h.fermRepo.Save(&fermentation); err != nil {
// 		h.HandleError(w, err)
// 	} else {
// 		h.Encode(w, &fermentation)
// 	}
// }

// func (h *FermentationHandler) handlePostTemperatureChange(w http.ResponseWriter, r *http.Request) {
// 	change, err := parseTemperatureChange(r)
// 	if err != nil {
// 		h.HandleError(w, err)
// 		return
// 	}
// 	if err := h.changeRepo.Save(&change); err != nil {
// 		h.HandleError(w, err)
// 	} else {
// 		h.Encode(w, &change)
// 	}
// }

// func (h *FermentationHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "" {
// 		if id, err := strconv.ParseUint(r.URL.Path, 10, 64); err != nil {
// 			h.HandleError(w, app.ErrBadRequest)
// 		} else {
// 			if err := h.fermRepo.Delete(id); err != nil {
// 				h.HandleError(w, err)
// 			}

// 			// ToDo: delete temperaturechanges

// 		}
// 		return
// 	}
// 	h.HandleError(w, app.ErrBadRequest)
// }

// func parseFermentation(r *http.Request) (internal.Fermentation, error) {
// 	var fermentation internal.Fermentation
// 	err := json.NewDecoder(r.Body).Decode(&fermentation)
// 	return fermentation, err
// }

// func parseTemperatureChange(r *http.Request) (internal.TemperatureChange, error) {
// 	var change internal.TemperatureChange
// 	err := json.NewDecoder(r.Body).Decode(&change)
// 	return change, err
// }
