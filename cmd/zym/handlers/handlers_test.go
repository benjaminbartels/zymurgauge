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
	batchID          = "KBTM3F9soO5TtbAx0A5mBZTAUsNZyg"
	repoErrMsg       = "could not %s repository"
	controllerErrMsg = "could not %s controller"
	respondErrMsg    = "problem responding to client"
	notFoundErrorMsg = "%s '%s' not found"
	parseErrorMsg    = "could not parse chamber"
)

var errSomeError = errors.New("some error")

func setupHandlerTest(query string, body io.Reader) (w *httptest.ResponseRecorder, r *http.Request,
	ctx context.Context) {
	w = httptest.NewRecorder()

	if query != "" {
		query = "?" + query
	}

	r = httptest.NewRequest("", "/"+query, body)
	v := web.CtxValues{
		Path: r.URL.Path,
		Now:  time.Now(),
	}
	ctx = web.InitContextValues(r.Context(), &v)

	return
}
