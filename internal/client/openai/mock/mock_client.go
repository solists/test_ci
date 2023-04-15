// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock_openai is a generated GoMock package.
package mock_openai

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	openai "github.com/sashabaranov/go-openai"
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
func (m *MockClient) GetQuery(ctx context.Context, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuery", ctx, messages)
	ret0, _ := ret[0].(*openai.ChatCompletionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuery indicates an expected call of GetQuery.
func (mr *MockClientMockRecorder) GetQuery(ctx, messages interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuery", reflect.TypeOf((*MockClient)(nil).GetQuery), ctx, messages)
}

// GetTranscription mocks base method.
func (m *MockClient) GetTranscription(ctx context.Context, filePath string) (*openai.AudioResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTranscription", ctx, filePath)
	ret0, _ := ret[0].(*openai.AudioResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTranscription indicates an expected call of GetTranscription.
func (mr *MockClientMockRecorder) GetTranscription(ctx, filePath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTranscription", reflect.TypeOf((*MockClient)(nil).GetTranscription), ctx, filePath)
}
