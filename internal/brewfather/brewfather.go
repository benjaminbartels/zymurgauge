package brewfather

import "context"

type RecipeRepo interface {
	GetRecipes(ctx context.Context) ([]Recipe, error)
	GetRecipe(ctx context.Context, id string) (*Recipe, error)
}

// TODO: Move this?
