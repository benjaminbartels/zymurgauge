package handlers

import (
	"net/http"
	"strconv"

	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

type API struct {
	http.Handler
	BeerHandler         *BeerHandler
	ChamberHandler      *ChamberHandler
	FermentationHandler *FermentationHandler
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var head string

	head, r.URL.Path = shiftPath(r.URL.Path)

	switch head {

	case "v1":

		head, r.URL.Path = shiftPath(r.URL.Path)

		switch head {

		case "chambers":
			a.ChamberHandler.ServeHTTP(w, r)
		case "beers":
			a.BeerHandler.ServeHTTP(w, r)
		case "fermentations":
			a.FermentationHandler.ServeHTTP(w, r)

		default:
			web.HandleError(w, web.ErrNotFound)
		}

	default:
		web.HandleError(w, web.ErrNotFound)

	}

}

func parseID(r *http.Request) (uint64, error) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	return strconv.ParseUint(head, 10, 64)
}
