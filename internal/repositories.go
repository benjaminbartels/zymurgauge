package internal

import (
	"context"

	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/recipe"
)

type ChamberRepo interface {
	GetAll() ([]chamber.Chamber, error)
	Get(id string) (*chamber.Chamber, error)
	Save(c *chamber.Chamber) error
	Delete(id string) error
}

type ThermometerRepo interface {
	GetThermometerIDs() ([]string, error)
}

type RecipeRepo interface {
	GetRecipes(ctx context.Context) ([]recipe.Recipe, error)
	GetRecipe(ctx context.Context, id string) (*recipe.Recipe, error)
}
