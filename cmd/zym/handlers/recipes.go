package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/recipe"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type RecipesHandler struct {
	Repo recipe.Repo
}

func (h *RecipesHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request,
	p httprouter.Params) error {
	recipes, err := h.Repo.GetRecipes(ctx)
	if err != nil {
		return errors.Wrap(err, "could not get all recipes from repository")
	}

	if err = web.Respond(ctx, w, recipes, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}

func (h *RecipesHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	recipe, err := h.Repo.GetRecipe(ctx, id)
	if err != nil {
		if errors.Is(err, brewfather.ErrNotFound) {
			return web.NewRequestError(fmt.Sprintf("recipe '%s' not found", id), http.StatusNotFound)
		}

		return errors.Wrap(err, "could not get recipe from repository")
	}

	if recipe == nil {
		return web.NewRequestError(fmt.Sprintf("recipe '%s' not found", id), http.StatusNotFound)
	}

	if err = web.Respond(ctx, w, recipe, http.StatusOK); err != nil {
		return errors.Wrap(err, "problem responding to client")
	}

	return nil
}
