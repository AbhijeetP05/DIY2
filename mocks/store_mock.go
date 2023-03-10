// Code generated by MockGen. DO NOT EDIT.
// Source: go-mux/services (interfaces: IStores)

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIStores is a mock of IStores interface.
type MockIStores struct {
	ctrl     *gomock.Controller
	recorder *MockIStoresMockRecorder
}

// MockIStoresMockRecorder is the mock recorder for MockIStores.
type MockIStoresMockRecorder struct {
	mock *MockIStores
}

// NewMockIStores creates a new mock instance.
func NewMockIStores(ctrl *gomock.Controller) *MockIStores {
	mock := &MockIStores{ctrl: ctrl}
	mock.recorder = &MockIStoresMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIStores) EXPECT() *MockIStoresMockRecorder {
	return m.recorder
}

// AddProducts mocks base method.
func (m *MockIStores) AddProducts(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddProducts", arg0, arg1)
}

// AddProducts indicates an expected call of AddProducts.
func (mr *MockIStoresMockRecorder) AddProducts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProducts", reflect.TypeOf((*MockIStores)(nil).AddProducts), arg0, arg1)
}

// GetProducts mocks base method.
func (m *MockIStores) GetProducts(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetProducts", arg0, arg1)
}

// GetProducts indicates an expected call of GetProducts.
func (mr *MockIStoresMockRecorder) GetProducts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducts", reflect.TypeOf((*MockIStores)(nil).GetProducts), arg0, arg1)
}
