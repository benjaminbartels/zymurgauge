package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/benjaminbartels/zymurgauge/internal/auth"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
)

const partsLength = 2

func Authorize(secret string) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
			authHeader := r.Header["Authorization"]

			if len(authHeader) != 1 {
				return web.NewRequestError("authorization header invalid", http.StatusBadRequest)
			}

			parts := strings.Split(authHeader[0], "Bearer")
			if len(parts) != partsLength {
				return web.NewRequestError("authorization header invalid", http.StatusBadRequest)
			}

			token := strings.TrimSpace(parts[1])
			if len(token) < 1 {
				return web.NewRequestError("authorization token invalid", http.StatusBadRequest)
			}

			if _, err := auth.IsAuthorized(secret, token); err != nil {
				return web.NewRequestError("access denied", http.StatusUnauthorized)
			}

			return handler(ctx, w, r, p)
		}

		return h
	}

	return m
}
