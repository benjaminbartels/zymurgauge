package handlers

import (
	"fmt"
	"net/http"

	"github.com/benjaminbartels/zymurgauge/cmd/fermmon/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
)

const chambersPath = "/recipes"
const recipesPath = "/recipes"

func NewAPI(chamberRepo *storage.ChamberRepo, service *brewfather.Service) http.Handler { // TODO: Better naming
	chambersHandler := &Chambers{
		repo: chamberRepo,
	}

	recipesHandler := &Recipes{
		service: service,
	}

	app := web.NewApp()
	app.Handle(http.MethodGet, chambersPath, chambersHandler.GetAll, false)
	app.Handle(http.MethodGet, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Get, false)
	app.Handle(http.MethodPost, chambersPath, chambersHandler.Save, false)
	app.Handle(http.MethodDelete, fmt.Sprintf("%s/:id", chambersPath), chambersHandler.Delete, false)
	app.Handle(http.MethodGet, recipesPath, recipesHandler.GetAll, false)
	app.Handle(http.MethodGet, fmt.Sprintf("%s/:id", recipesPath), recipesHandler.Get, false)

	return app
}
