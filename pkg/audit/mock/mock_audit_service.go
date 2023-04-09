// Code generated by MockGen. DO NOT EDIT.
// Source: audit_service.go

// Package mock_audit is a generated GoMock package.
package mock_audit

import (
	audit "mymod/pkg/audit"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Log mocks base method.
func (m *MockService) Log(log *audit.Log) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Log", log)
}

// Log indicates an expected call of Log.
func (mr *MockServiceMockRecorder) Log(log interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockService)(nil).Log), log)
}
