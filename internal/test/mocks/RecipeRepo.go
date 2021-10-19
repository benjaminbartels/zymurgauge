// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	brewfather "github.com/benjaminbartels/zymurgauge/internal/brewfather"

	mock "github.com/stretchr/testify/mock"
)

// RecipeRepo is an autogenerated mock type for the RecipeRepo type
type RecipeRepo struct {
	mock.Mock
}

// GetRecipe provides a mock function with given fields: ctx, id
func (_m *RecipeRepo) GetRecipe(ctx context.Context, id string) (*brewfather.Recipe, error) {
	ret := _m.Called(ctx, id)

	var r0 *brewfather.Recipe
	if rf, ok := ret.Get(0).(func(context.Context, string) *brewfather.Recipe); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*brewfather.Recipe)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRecipes provides a mock function with given fields: ctx
func (_m *RecipeRepo) GetRecipes(ctx context.Context) ([]brewfather.Recipe, error) {
	ret := _m.Called(ctx)

	var r0 []brewfather.Recipe
	if rf, ok := ret.Get(0).(func(context.Context) []brewfather.Recipe); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]brewfather.Recipe)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
