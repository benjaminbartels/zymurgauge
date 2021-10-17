package handlers_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

const (
	chamberID        = "96f58a65-03c0-49f3-83ca-ab751bbf3768"
	recipeID         = "0789f223-b5fb-49c1-a8a8-111adff88b82"
	repoErrMsg       = "could not %s repository"
	respondErrMsg    = "problem responding to client"
	notFoundErrorMsg = "%s '%s' not found"
	parseErrorMsg    = "could not parse chamber"
)

var errDeadDatabase = errors.New("database is dead")

func setupHandlerTest(body io.Reader) (w *httptest.ResponseRecorder, r *http.Request, ctx context.Context) {
	w = httptest.NewRecorder()
	r = httptest.NewRequest("", "/", body)
	v := web.CtxValues{
		Path: r.URL.Path,
		Now:  time.Now(),
	}
	ctx = web.InitContextValues(r.Context(), &v)

	return
}
