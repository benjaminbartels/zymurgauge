package middleware

import (
	"context"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
)

func Cors() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
			w.Header().Set("Access-Control-Allow-Origin", "*")

			return handler(ctx, w, r, p)
		}

		return h
	}

	return m
}
