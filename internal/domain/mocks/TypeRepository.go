// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kamkali/go-timeline/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// TypeRepository is an autogenerated mock type for the TypeRepository type
type TypeRepository struct {
	mock.Mock
}

// CreateType provides a mock function with given fields: ctx, t
func (_m *TypeRepository) CreateType(ctx context.Context, t *domain.Type) (domain.Type, error) {
	ret := _m.Called(ctx, t)

	var r0 domain.Type
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Type) domain.Type); ok {
		r0 = rf(ctx, t)
	} else {
		r0 = ret.Get(0).(domain.Type)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.Type) error); ok {
		r1 = rf(ctx, t)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTypeRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewTypeRepository creates a new instance of TypeRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTypeRepository(t mockConstructorTestingTNewTypeRepository) *TypeRepository {
	mock := &TypeRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}