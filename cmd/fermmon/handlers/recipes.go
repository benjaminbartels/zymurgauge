package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/benjaminbartels/zymurgauge/cmd/fermmon/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/julienschmidt/httprouter"
)

type Recipes struct {
	Service *brewfather.Service
}

func (h *Recipes) GetRecipes(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	recipes, err := h.Service.GetRecipes(ctx)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, recipes, http.StatusOK)
}

func (h *Recipes) GetRecipe(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id := p.ByName("id")

	recipe, err := h.Service.GetRecipe(ctx, id)
	if err != nil {
		fmt.Println("Could not get Recipes:", err)
		os.Exit(1)
	}

	if recipe == nil {
		return web.ErrNotFound
	}

	return web.Respond(ctx, w, recipe, http.StatusOK)
}
