package recipe

import "context"

type Repo interface {
	GetRecipes(ctx context.Context) ([]Recipe, error)
	GetRecipe(ctx context.Context, id string) (*Recipe, error)
}
