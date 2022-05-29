package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/benjaminbartels/zymurgauge/internal/auth"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/settings"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	expiresIn = 7 * 24 * time.Hour
)

type AuthHandler struct {
	SettingsRepo settings.Repo
	Logger       *logrus.Logger
}

func (h *AuthHandler) Login(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	creds, err := parseCredentials(r)
	if err != nil {
		return errors.Wrap(err, "could not parse credentials")
	}

	s, err := h.SettingsRepo.Get()
	if err != nil {
		return errors.Wrap(err, "could not get settings")
	}

	if creds.Username != s.Username {
		return web.NewRequestError("incorrect username", http.StatusUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(creds.Password)); err != nil {
		return web.NewRequestError("incorrect password", http.StatusUnauthorized)
	}

	token, err := auth.CreateToken(s.AuthSecret, creds.Username, expiresIn)
	if err != nil {
		return errors.Wrap(err, "could not create token")
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	if err := web.Respond(ctx, w, response, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func parseCredentials(r *http.Request) (auth.Credentials, error) {
	var user auth.Credentials
	err := json.NewDecoder(r.Body).Decode(&user)

	return user, errors.Wrap(err, "could not decode user from request body")
}
