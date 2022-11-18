// Code generated by MockGen. DO NOT EDIT.
// Source: internal/messages/incoming_msg.go

// Package mock_messages is a generated GoMock package.
package mock_messages

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMessageSender is a mock of MessageSender interface.
type MockMessageSender struct {
	ctrl     *gomock.Controller
	recorder *MockMessageSenderMockRecorder
}

// MockMessageSenderMockRecorder is the mock recorder for MockMessageSender.
type MockMessageSenderMockRecorder struct {
	mock *MockMessageSender
}

// NewMockMessageSender creates a new mock instance.
func NewMockMessageSender(ctrl *gomock.Controller) *MockMessageSender {
	mock := &MockMessageSender{ctrl: ctrl}
	mock.recorder = &MockMessageSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageSender) EXPECT() *MockMessageSenderMockRecorder {
	return m.recorder
}

// SendCurrencyKeyboard mocks base method.
func (m *MockMessageSender) SendCurrencyKeyboard(ctx context.Context, userID int64, text string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCurrencyKeyboard", ctx, userID, text)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCurrencyKeyboard indicates an expected call of SendCurrencyKeyboard.
func (mr *MockMessageSenderMockRecorder) SendCurrencyKeyboard(ctx, userID, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCurrencyKeyboard", reflect.TypeOf((*MockMessageSender)(nil).SendCurrencyKeyboard), ctx, userID, text)
}

// SendMessage mocks base method.
func (m *MockMessageSender) SendMessage(ctx context.Context, text string, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", ctx, text, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockMessageSenderMockRecorder) SendMessage(ctx, text, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockMessageSender)(nil).SendMessage), ctx, text, userID)
}

// SendRangeKeyboard mocks base method.
func (m *MockMessageSender) SendRangeKeyboard(ctx context.Context, userID int64, text string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRangeKeyboard", ctx, userID, text)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendRangeKeyboard indicates an expected call of SendRangeKeyboard.
func (mr *MockMessageSenderMockRecorder) SendRangeKeyboard(ctx, userID, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRangeKeyboard", reflect.TypeOf((*MockMessageSender)(nil).SendRangeKeyboard), ctx, userID, text)
}
