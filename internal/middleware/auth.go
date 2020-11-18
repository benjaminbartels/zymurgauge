package middleware

import (
	"context"
	"net/http"

	"github.com/auth0-community/go-auth0"
	"github.com/benjaminbartels/zymurgauge/internal/platform/log"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/pkg/errors"
	jose "gopkg.in/square/go-jose.v2"
)

const (
	audience = "zymurgauge.com/api"
	issuer   = "https://zymurgauge.auth0.com/"
)

// Authorizer check validates that the authorization header of the request.
type Authorizer struct {
	clientSecret string
	logger       log.Logger
}

// NewAuthorizer creates a new Authorizer.
func NewAuthorizer(clientSecret string, logger log.Logger) *Authorizer {
	return &Authorizer{
		clientSecret: clientSecret,
		logger:       logger,
	}
}

// Authorize validates that token contained authorization header.  If the token is invalid a 401 unauthorized status is
// returned in the response.
func (a *Authorizer) Authorize(next web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		if r.Method == "OPTIONS" { // ToDo: And more to circumvent auth check (login?, get-token?)
			return next(ctx, w, r)
		}

		secretProvider := auth0.NewKeyProvider([]byte(a.clientSecret))

		configuration := auth0.NewConfiguration(secretProvider, []string{audience}, issuer, jose.HS256)
		validator := auth0.NewValidator(configuration, auth0.RequestTokenExtractorFunc(auth0.FromHeader))

		a.logger.Println(r.Header["Authorization"])

		_, err := validator.ValidateRequest(r)
		if err != nil {
			return errors.WithMessage(web.ErrUnauthorized, err.Error())
		}

		return next(ctx, w, r)
	}

	return h
}
