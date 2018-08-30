package middleware

import (
	"context"
	"net/http"

	"github.com/auth0-community/auth0"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v2"
)

type Authorizer struct {
	clientSecret string
	logger       log.Logger
}

func NewAuthorizer(clientSecret string, logger log.Logger) *Authorizer {
	return &Authorizer{
		clientSecret: clientSecret,
		logger:       logger,
	}
}

func (a *Authorizer) Authorize(next web.Handler) web.Handler {

	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

		if r.Method == "OPTIONS" { // ToDo: And more to circumvent auth check (login?, get-token?)
			return next(ctx, w, r)
		}

		secretProvider := auth0.NewKeyProvider([]byte(a.clientSecret))
		audience := []string{"https://www.zymurgauge.com/api"}

		configuration := auth0.NewConfiguration(secretProvider, audience, "https://zymurgauge.auth0.com/", jose.HS256)
		validator := auth0.NewValidator(configuration, auth0.RequestTokenExtractorFunc(auth0.FromHeader))

		_, err := validator.ValidateRequest(r)

		if err != nil {
			return errors.WithMessage(web.ErrUnauthorized, err.Error())
		}

		return next(ctx, w, r)
	}

	return h

}
