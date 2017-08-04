package handlers

import (
	"net/http"
	"path"
	"strconv"
	"strings"

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

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func parseID(r *http.Request) (uint64, error) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	return strconv.ParseUint(head, 10, 64)
}
