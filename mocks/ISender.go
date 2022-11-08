// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ISender is an autogenerated mock type for the ISender type
type ISender struct {
	mock.Mock
}

// SendMessage provides a mock function with given fields: _a0, destination, message
func (_m *ISender) SendMessage(_a0 context.Context, destination string, message []byte) error {
	ret := _m.Called(_a0, destination, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte) error); ok {
		r0 = rf(_a0, destination, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewISender interface {
	mock.TestingT
	Cleanup(func())
}

// NewISender creates a new instance of ISender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewISender(t mockConstructorTestingTNewISender) *ISender {
	mock := &ISender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}