// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	bluetooth "github.com/benjaminbartels/zymurgauge/internal/platform/bluetooth"

	linux "github.com/go-ble/ble/linux"

	mock "github.com/stretchr/testify/mock"
)

// Scanner is an autogenerated mock type for the Scanner type
type Scanner struct {
	mock.Mock
}

// NewDevice provides a mock function with given fields:
func (_m *Scanner) NewDevice() (*linux.Device, error) {
	ret := _m.Called()

	var r0 *linux.Device
	if rf, ok := ret.Get(0).(func() *linux.Device); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*linux.Device)
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

// Scan provides a mock function with given fields: ctx, h, f
func (_m *Scanner) Scan(ctx context.Context, h func(bluetooth.Advertisement), f func(bluetooth.Advertisement) bool) error {
	ret := _m.Called(ctx, h, f)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(bluetooth.Advertisement), func(bluetooth.Advertisement) bool) error); ok {
		r0 = rf(ctx, h, f)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetDefaultDevice provides a mock function with given fields: device
func (_m *Scanner) SetDefaultDevice(device bluetooth.Device) {
	_m.Called(device)
}

// WithSigHandler provides a mock function with given fields: ctx, cancel
func (_m *Scanner) WithSigHandler(ctx context.Context, cancel func()) context.Context {
	ret := _m.Called(ctx, cancel)

	var r0 context.Context
	if rf, ok := ret.Get(0).(func(context.Context, func()) context.Context); ok {
		r0 = rf(ctx, cancel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}
