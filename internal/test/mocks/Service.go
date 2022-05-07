// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	brewfather "github.com/benjaminbartels/zymurgauge/internal/brewfather"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, id
func (_m *Service) GetBatchDetail(ctx context.Context, id string) (*brewfather.BatchDetail, error) {
	ret := _m.Called(ctx, id)

	var r0 *brewfather.BatchDetail
	if rf, ok := ret.Get(0).(func(context.Context, string) *brewfather.BatchDetail); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*brewfather.BatchDetail)
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

// GetAll provides a mock function with given fields: ctx
func (_m *Service) GetAllBatchSummaries(ctx context.Context) ([]brewfather.BatchSummary, error) {
	ret := _m.Called(ctx)

	var r0 []brewfather.BatchSummary
	if rf, ok := ret.Get(0).(func(context.Context) []brewfather.BatchSummary); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]brewfather.BatchSummary)
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

// Log provides a mock function with given fields: ctx, log
func (_m *Service) Log(ctx context.Context, log brewfather.LogEntry) error {
	ret := _m.Called(ctx, log)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, brewfather.LogEntry) error); ok {
		r0 = rf(ctx, log)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
