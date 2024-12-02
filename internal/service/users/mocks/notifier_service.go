// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	net "net"

	mock "github.com/stretchr/testify/mock"
)

// NotifierService is an autogenerated mock type for the NotifierService type
type NotifierService struct {
	mock.Mock
}

// SendSuspiciousActivityMail provides a mock function with given fields: email, newIP
func (_m *NotifierService) SendSuspiciousActivityMail(email string, newIP net.IP) error {
	ret := _m.Called(email, newIP)

	if len(ret) == 0 {
		panic("no return value specified for SendSuspiciousActivityMail")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, net.IP) error); ok {
		r0 = rf(email, newIP)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewNotifierService creates a new instance of NotifierService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNotifierService(t interface {
	mock.TestingT
	Cleanup(func())
}) *NotifierService {
	mock := &NotifierService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
