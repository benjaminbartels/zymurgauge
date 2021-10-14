// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	storage "github.com/benjaminbartels/zymurgauge/internal/storage"
	mock "github.com/stretchr/testify/mock"
)

// ChamberRepo is an autogenerated mock type for the ChamberRepo type
type ChamberRepo struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *ChamberRepo) Delete(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: id
func (_m *ChamberRepo) Get(id string) (*storage.Chamber, error) {
	ret := _m.Called(id)

	var r0 *storage.Chamber
	if rf, ok := ret.Get(0).(func(string) *storage.Chamber); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.Chamber)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *ChamberRepo) GetAll() ([]storage.Chamber, error) {
	ret := _m.Called()

	var r0 []storage.Chamber
	if rf, ok := ret.Get(0).(func() []storage.Chamber); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]storage.Chamber)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: c
func (_m *ChamberRepo) Save(c *storage.Chamber) error {
	ret := _m.Called(c)

	var r0 error
	if rf, ok := ret.Get(0).(func(*storage.Chamber) error); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}