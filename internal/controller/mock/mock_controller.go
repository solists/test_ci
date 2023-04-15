// Code generated by MockGen. DO NOT EDIT.
// Source: controller.go

// Package mock_controller is a generated GoMock package.
package mock_controller

import (
	context "context"
	openai "mymod/internal/models/openai"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIController is a mock of IController interface.
type MockIController struct {
	ctrl     *gomock.Controller
	recorder *MockIControllerMockRecorder
}

// MockIControllerMockRecorder is the mock recorder for MockIController.
type MockIControllerMockRecorder struct {
	mock *MockIController
}

// NewMockIController creates a new mock instance.
func NewMockIController(ctrl *gomock.Controller) *MockIController {
	mock := &MockIController{ctrl: ctrl}
	mock.recorder = &MockIControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIController) EXPECT() *MockIControllerMockRecorder {
	return m.recorder
}

// GetQuery mocks base method.
func (m *MockIController) GetQuery(ctx context.Context, req *openai.GetQueryRequest) (*openai.GetQueryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuery", ctx, req)
	ret0, _ := ret[0].(*openai.GetQueryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuery indicates an expected call of GetQuery.
func (mr *MockIControllerMockRecorder) GetQuery(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuery", reflect.TypeOf((*MockIController)(nil).GetQuery), ctx, req)
}

// GetTranscription mocks base method.
func (m *MockIController) GetTranscription(ctx context.Context, req *openai.GetTranscriptionRequest) (*openai.GetTranscriptionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTranscription", ctx, req)
	ret0, _ := ret[0].(*openai.GetTranscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTranscription indicates an expected call of GetTranscription.
func (mr *MockIControllerMockRecorder) GetTranscription(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTranscription", reflect.TypeOf((*MockIController)(nil).GetTranscription), ctx, req)
}
