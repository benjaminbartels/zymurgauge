// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Metrics is an autogenerated mock type for the Metrics type
type Metrics struct {
	mock.Mock
}

// Count provides a mock function with given fields: bucket, n
func (_m *Metrics) Count(bucket string, n interface{}) {
	_m.Called(bucket, n)
}

// Gauge provides a mock function with given fields: bucket, value
func (_m *Metrics) Gauge(bucket string, value interface{}) {
	_m.Called(bucket, value)
}

// Histogram provides a mock function with given fields: bucket, value
func (_m *Metrics) Histogram(bucket string, value interface{}) {
	_m.Called(bucket, value)
}

// Increment provides a mock function with given fields: bucket
func (_m *Metrics) Increment(bucket string) {
	_m.Called(bucket)
}

// Timing provides a mock function with given fields: bucket, value
func (_m *Metrics) Timing(bucket string, value interface{}) {
	_m.Called(bucket, value)
}
