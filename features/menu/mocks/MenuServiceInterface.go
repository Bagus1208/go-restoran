// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	model "restoran/features/menu/model"

	mock "github.com/stretchr/testify/mock"

	multipart "mime/multipart"
)

// MenuServiceInterface is an autogenerated mock type for the MenuServiceInterface type
type MenuServiceInterface struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *MenuServiceInterface) Delete(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: pagination
func (_m *MenuServiceInterface) GetAll(pagination model.QueryParam) ([]model.MenuResponse, error) {
	ret := _m.Called(pagination)

	var r0 []model.MenuResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(model.QueryParam) ([]model.MenuResponse, error)); ok {
		return rf(pagination)
	}
	if rf, ok := ret.Get(0).(func(model.QueryParam) []model.MenuResponse); ok {
		r0 = rf(pagination)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.MenuResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(model.QueryParam) error); ok {
		r1 = rf(pagination)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: name
func (_m *MenuServiceInterface) GetByName(name string) (model.MenuResponse, error) {
	ret := _m.Called(name)

	var r0 model.MenuResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (model.MenuResponse, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) model.MenuResponse); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(model.MenuResponse)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCategory provides a mock function with given fields: queryParam
func (_m *MenuServiceInterface) GetCategory(queryParam model.QueryParam) ([]model.MenuResponse, error) {
	ret := _m.Called(queryParam)

	var r0 []model.MenuResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(model.QueryParam) ([]model.MenuResponse, error)); ok {
		return rf(queryParam)
	}
	if rf, ok := ret.Get(0).(func(model.QueryParam) []model.MenuResponse); ok {
		r0 = rf(queryParam)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.MenuResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(model.QueryParam) error); ok {
		r1 = rf(queryParam)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFavorite provides a mock function with given fields:
func (_m *MenuServiceInterface) GetFavorite() ([]model.Favorite, error) {
	ret := _m.Called()

	var r0 []model.Favorite
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]model.Favorite, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []model.Favorite); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Favorite)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: fileHeader, newData
func (_m *MenuServiceInterface) Insert(fileHeader *multipart.FileHeader, newData model.MenuInput) (*model.MenuResponse, error) {
	ret := _m.Called(fileHeader, newData)

	var r0 *model.MenuResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*multipart.FileHeader, model.MenuInput) (*model.MenuResponse, error)); ok {
		return rf(fileHeader, newData)
	}
	if rf, ok := ret.Get(0).(func(*multipart.FileHeader, model.MenuInput) *model.MenuResponse); ok {
		r0 = rf(fileHeader, newData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MenuResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*multipart.FileHeader, model.MenuInput) error); ok {
		r1 = rf(fileHeader, newData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: id, fileHeader, updateData
func (_m *MenuServiceInterface) Update(id int, fileHeader *multipart.FileHeader, updateData model.MenuInput) (*model.MenuResponse, error) {
	ret := _m.Called(id, fileHeader, updateData)

	var r0 *model.MenuResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(int, *multipart.FileHeader, model.MenuInput) (*model.MenuResponse, error)); ok {
		return rf(id, fileHeader, updateData)
	}
	if rf, ok := ret.Get(0).(func(int, *multipart.FileHeader, model.MenuInput) *model.MenuResponse); ok {
		r0 = rf(id, fileHeader, updateData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MenuResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(int, *multipart.FileHeader, model.MenuInput) error); ok {
		r1 = rf(id, fileHeader, updateData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMenuServiceInterface creates a new instance of MenuServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMenuServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MenuServiceInterface {
	mock := &MenuServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}