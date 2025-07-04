// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	models "github.com/vaidik-bajpai/medibridge/internal/models"
)

// VitalsStorer is an autogenerated mock type for the VitalsStorer type
type VitalsStorer struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, req
func (_m *VitalsStorer) Create(ctx context.Context, req *models.CreateVitalReq) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.CreateVitalReq) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, pID
func (_m *VitalsStorer) Delete(ctx context.Context, pID string) error {
	ret := _m.Called(ctx, pID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, pID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, req
func (_m *VitalsStorer) Update(ctx context.Context, req *models.UpdateVitalReq) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.UpdateVitalReq) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewVitalsStorer creates a new instance of VitalsStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewVitalsStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *VitalsStorer {
	mock := &VitalsStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
