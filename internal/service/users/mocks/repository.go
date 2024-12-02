// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	users "github.com/SpaceSlow/test-task-backend-junior-medods/internal/domain/users"

	uuid "github.com/google/uuid"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// CreateRefreshToken provides a mock function with given fields: userGUID, refresh
func (_m *Repository) CreateRefreshToken(userGUID uuid.UUID, refresh *users.RefreshToken) error {
	ret := _m.Called(userGUID, refresh)

	if len(ret) == 0 {
		panic("no return value specified for CreateRefreshToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, *users.RefreshToken) error); ok {
		r0 = rf(userGUID, refresh)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchEmailByUUID provides a mock function with given fields: userGUID
func (_m *Repository) FetchEmailByUUID(userGUID uuid.UUID) (string, error) {
	ret := _m.Called(userGUID)

	if len(ret) == 0 {
		panic("no return value specified for FetchEmailByUUID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (string, error)); ok {
		return rf(userGUID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) string); ok {
		r0 = rf(userGUID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(userGUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchUserByEmail provides a mock function with given fields: email
func (_m *Repository) FetchUserByEmail(email string) (*users.User, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for FetchUserByEmail")
	}

	var r0 *users.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*users.User, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *users.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
