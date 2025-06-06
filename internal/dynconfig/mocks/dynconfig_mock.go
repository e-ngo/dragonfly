// Code generated by MockGen. DO NOT EDIT.
// Source: dynconfig.go
//
// Generated by this command:
//
//	mockgen -destination mocks/dynconfig_mock.go -source dynconfig.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDynconfig is a mock of Dynconfig interface.
type MockDynconfig[T any] struct {
	ctrl     *gomock.Controller
	recorder *MockDynconfigMockRecorder[T]
	isgomock struct{}
}

// MockDynconfigMockRecorder is the mock recorder for MockDynconfig.
type MockDynconfigMockRecorder[T any] struct {
	mock *MockDynconfig[T]
}

// NewMockDynconfig creates a new mock instance.
func NewMockDynconfig[T any](ctrl *gomock.Controller) *MockDynconfig[T] {
	mock := &MockDynconfig[T]{ctrl: ctrl}
	mock.recorder = &MockDynconfigMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDynconfig[T]) EXPECT() *MockDynconfigMockRecorder[T] {
	return m.recorder
}

// Get mocks base method.
func (m *MockDynconfig[T]) Get() (*T, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(*T)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDynconfigMockRecorder[T]) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDynconfig[T])(nil).Get))
}

// Refresh mocks base method.
func (m *MockDynconfig[T]) Refresh() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh")
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh.
func (mr *MockDynconfigMockRecorder[T]) Refresh() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockDynconfig[T])(nil).Refresh))
}
