package handlers

import (
	"fmt"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/cmd/fermmon/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

const recipesPath = "/recipes"

func NewAPI(service *brewfather.Service) http.Handler { // TODO: Better naming
	recipesHandler := &Recipes{
		Service: service,
	}

	app := web.NewApp()
	app.Handle(http.MethodGet, recipesPath, recipesHandler.GetRecipes, false)
	app.Handle(http.MethodGet, fmt.Sprintf("%s/:id", recipesPath), recipesHandler.GetRecipe, false)

	return app
}
