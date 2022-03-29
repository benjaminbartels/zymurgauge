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
)

const (
	expiresIn = 7 * 24 * time.Hour
)

type AuthHandler struct {
	UserRepo     auth.UserRepo
	SettingsRepo settings.Repo
	Logger       *logrus.Logger // TODO: use logrus.Entry
}

func (h *AuthHandler) Login(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	user, err := parseUser(r)
	if err != nil {
		return errors.Wrap(err, "could not parse credentials")
	}

	u, err := h.UserRepo.Get()
	if err != nil {
		return errors.Wrap(err, "could not get settings")
	}

	if user.Username != u.Username || user.Password != u.Password {
		return web.NewRequestError("incorrect username and/or password", http.StatusUnauthorized)
	}

	s, err := h.SettingsRepo.Get()
	if err != nil {
		return errors.Wrap(err, "could not get settings")
	}

	token, err := auth.CreateToken(s.AuthSecret, user, expiresIn)
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

func parseUser(r *http.Request) (auth.User, error) {
	var user auth.User
	err := json.NewDecoder(r.Body).Decode(&user)

	return user, errors.Wrap(err, "could not decode user from request body")
}
