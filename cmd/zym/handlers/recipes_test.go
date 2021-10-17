package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/test/mocks"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetAllRecipes(t *testing.T) {
	t.Parallel()
	t.Run("getAllRecipes", getAllRecipes)
	t.Run("getAllRecipesEmpty", getAllRecipesEmpty)
	t.Run("getAllRecipesRepoError", getAllRecipesRepoError)
	t.Run("getAllRecipesRespondError", getAllRecipesRespondError)
}

func getAllRecipes(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)

	expected := []brewfather.Recipe{
		{ID: recipeID},
		{ID: "f4ce0e05-1ada-42b8-8fc4-fb3482525d0d"},
	}
	repoMock := &mocks.RecipeRepo{}
	repoMock.On("GetRecipes", ctx).Return(expected, nil)

	handler := &handlers.RecipesHandler{Repo: repoMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []brewfather.Recipe{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func getAllRecipesEmpty(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)

	expected := []brewfather.Recipe{}
	repoMock := &mocks.RecipeRepo{}
	repoMock.On("GetRecipes", ctx).Return(expected, nil)

	handler := &handlers.RecipesHandler{Repo: repoMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	result := []brewfather.Recipe{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func getAllRecipesRepoError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)

	repoMock := &mocks.RecipeRepo{}
	repoMock.On("GetRecipes", ctx).Return([]brewfather.Recipe{}, errDeadDatabase)

	handler := &handlers.RecipesHandler{Repo: repoMock}
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	// TODO: Waiting on PR for ErrorContains(): https://github.com/stretchr/testify/pull/1022
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "get all recipes from"))
}

func getAllRecipesRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest(nil)
	ctx := context.Background()

	repoMock := &mocks.RecipeRepo{}
	repoMock.On("GetRecipes", ctx).Return([]brewfather.Recipe{}, nil)

	handler := &handlers.RecipesHandler{Repo: repoMock}
	// use new ctx to force error
	err := handler.GetAll(ctx, w, r, httprouter.Params{})
	assert.Contains(t, err.Error(), respondErrMsg)
}

//nolint: paralleltest // False positives with r.Run not in a loop
func TestGetRecipe(t *testing.T) {
	t.Parallel()
	t.Run("getRecipeFound", getRecipeFound)
	t.Run("getRecipeNotFound", getRecipeNotFound)
	t.Run("getRecipeNotFoundError", getRecipeNotFoundError)
	t.Run("getRecipeRepoError", getRecipeRepoError)
	t.Run("getRecipeRespondError", getRecipeRespondError)
}

func getRecipeFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)

	expected := brewfather.Recipe{ID: recipeID}
	repoMock := &mocks.RecipeRepo{}
	repoMock.On("GetRecipe", ctx, recipeID).Return(&expected, nil)

	handler := &handlers.RecipesHandler{Repo: repoMock}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: recipeID}})
	assert.NoError(t, err)

	resp := w.Result()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	recipe := brewfather.Recipe{}
	err = json.Unmarshal(bodyBytes, &recipe)
	assert.NoError(t, err)
	assert.Equal(t, expected, recipe)
}

func getRecipeNotFound(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)

	var expected *brewfather.Recipe

	repoMock := &mocks.RecipeRepo{}
	repoMock.On("GetRecipe", ctx, recipeID).Return(expected, nil)

	handler := &handlers.RecipesHandler{Repo: repoMock}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: recipeID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "recipe", recipeID))
}

func getRecipeNotFoundError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)

	var expected *brewfather.Recipe

	repoMock := &mocks.RecipeRepo{}
	repoMock.On("GetRecipe", ctx, recipeID).Return(expected, brewfather.ErrNotFound)

	handler := &handlers.RecipesHandler{Repo: repoMock}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: recipeID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(notFoundErrorMsg, "recipe", recipeID))
}

func getRecipeRepoError(t *testing.T) {
	t.Parallel()

	w, r, ctx := setupHandlerTest(nil)

	var expected *brewfather.Recipe

	repoMock := &mocks.RecipeRepo{}
	repoMock.On("GetRecipe", ctx, recipeID).Return(expected, errDeadDatabase)

	handler := &handlers.RecipesHandler{Repo: repoMock}
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: recipeID}})
	assert.Contains(t, err.Error(), fmt.Sprintf(repoErrMsg, "get recipe from"))
}

func getRecipeRespondError(t *testing.T) {
	t.Parallel()

	w, r, _ := setupHandlerTest(nil)
	ctx := context.Background()

	expected := brewfather.Recipe{ID: recipeID}
	repoMock := &mocks.RecipeRepo{}
	repoMock.On("GetRecipe", ctx, recipeID).Return(&expected, nil)

	handler := &handlers.RecipesHandler{Repo: repoMock}
	// use new ctx to force error
	err := handler.Get(ctx, w, r, httprouter.Params{httprouter.Param{Key: "id", Value: recipeID}})
	assert.Contains(t, err.Error(), respondErrMsg)
}
