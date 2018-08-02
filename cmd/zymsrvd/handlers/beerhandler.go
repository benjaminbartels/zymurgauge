package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

// BeerHandler is the http handler for API calls to manage Beers
type BeerHandler struct {
	repo *database.BeerRepo
}

// NewBeerHandler instantiates a BeerHandler
func NewBeerHandler(repo *database.BeerRepo) *BeerHandler {
	return &BeerHandler{
		repo: repo,
	}
}

func (h *BeerHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request, params web.Params) error {
	if beers, err := h.repo.GetAll(); err != nil {
		return err
	} else {
		return web.Respond(ctx, w, beers, http.StatusOK)
	}
}

func (h *BeerHandler) GetOne(ctx context.Context, w http.ResponseWriter, r *http.Request, params web.Params) error {
	id := params.ByName("id")

	beerID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err //ToDo: error InvalidID
	}

	if beer, err := h.repo.Get(beerID); err != nil {
		return err
	} else if beer == nil {
		return web.ErrNotFound
	} else {
		return web.Respond(ctx, w, beer, http.StatusOK)
	}
}

// func (h *BeerHandler) handlePost(w http.ResponseWriter, r *http.Request) {
// 	beer, err := parseBeer(r)
// 	if err != nil {
// 		h.HandleError(w, err)
// 		return
// 	}

// 	if err := h.repo.Save(&beer); err != nil {
// 		h.HandleError(w, err)
// 	} else {
// 		h.Encode(w, &beer)
// 	}
// }

// func (h *BeerHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "" {
// 		if id, err := strconv.ParseUint(r.URL.Path, 10, 64); err != nil {
// 			h.HandleError(w, app.ErrBadRequest)
// 		} else {
// 			if err := h.repo.Delete(id); err != nil {
// 				h.HandleError(w, err)
// 			}
// 		}
// 		return
// 	}
// 	h.HandleError(w, app.ErrBadRequest)
// }

// func parseBeer(r *http.Request) (internal.Beer, error) {
// 	var beer internal.Beer
// 	err := json.NewDecoder(r.Body).Decode(&beer)
// 	return beer, err
// }

// func (b *BeerHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	beers := []internal.Beer{internal.Beer{
// 		Name: "Golden Stout",
// 	}, internal.Beer{
// 		Name: "Barley Wine",
// 	}}

// 	// Do stuff and check for errors
// 	// if err = check(err); err != nil {
// 	// 	return errors.Wrap(err, "")
// 	// }
// 	err := web.Respond(ctx, w, beers, http.StatusOK)
// 	return err
// }
