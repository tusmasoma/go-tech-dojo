// Code generated by MockGen. DO NOT EDIT.
// Source: game.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	usecase "github.com/tusmasoma/go-tech-dojo/usecase"
)

// MockGameUseCase is a mock of GameUseCase interface.
type MockGameUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockGameUseCaseMockRecorder
}

// MockGameUseCaseMockRecorder is the mock recorder for MockGameUseCase.
type MockGameUseCaseMockRecorder struct {
	mock *MockGameUseCase
}

// NewMockGameUseCase creates a new mock instance.
func NewMockGameUseCase(ctrl *gomock.Controller) *MockGameUseCase {
	mock := &MockGameUseCase{ctrl: ctrl}
	mock.recorder = &MockGameUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGameUseCase) EXPECT() *MockGameUseCaseMockRecorder {
	return m.recorder
}

// DrawGacha mocks base method.
func (m *MockGameUseCase) DrawGacha(ctx context.Context, times int) ([]*usecase.GachaResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DrawGacha", ctx, times)
	ret0, _ := ret[0].([]*usecase.GachaResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DrawGacha indicates an expected call of DrawGacha.
func (mr *MockGameUseCaseMockRecorder) DrawGacha(ctx, times interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DrawGacha", reflect.TypeOf((*MockGameUseCase)(nil).DrawGacha), ctx, times)
}

// FinishGame mocks base method.
func (m *MockGameUseCase) FinishGame(ctx context.Context, scoreValue int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishGame", ctx, scoreValue)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FinishGame indicates an expected call of FinishGame.
func (mr *MockGameUseCaseMockRecorder) FinishGame(ctx, scoreValue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishGame", reflect.TypeOf((*MockGameUseCase)(nil).FinishGame), ctx, scoreValue)
}
