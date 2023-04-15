// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock_openai is a generated GoMock package.
package mock_openai

import (
	context "context"
	openai "mymod/internal/models/openai"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GetQuery mocks base method.
func (m *MockClient) GetQuery(ctx context.Context, req *openai.GetQueryRequest) (*openai.GetQueryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuery", ctx, req)
	ret0, _ := ret[0].(*openai.GetQueryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQueryOPENAI indicates an expected call of GetQueryOPENAI.
func (mr *MockClientMockRecorder) GetQueryOPENAI(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuery", reflect.TypeOf((*MockClient)(nil).GetQuery), ctx, req)
}
