// Code generated by MockGen. DO NOT EDIT.
// Source: ranking.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	model "github.com/tusmasoma/go-tech-dojo/domain/model"
)

// MockRankingUseCase is a mock of RankingUseCase interface.
type MockRankingUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockRankingUseCaseMockRecorder
}

// MockRankingUseCaseMockRecorder is the mock recorder for MockRankingUseCase.
type MockRankingUseCaseMockRecorder struct {
	mock *MockRankingUseCase
}

// NewMockRankingUseCase creates a new mock instance.
func NewMockRankingUseCase(ctrl *gomock.Controller) *MockRankingUseCase {
	mock := &MockRankingUseCase{ctrl: ctrl}
	mock.recorder = &MockRankingUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRankingUseCase) EXPECT() *MockRankingUseCaseMockRecorder {
	return m.recorder
}

// ListRankings mocks base method.
func (m *MockRankingUseCase) ListRankings(ctx context.Context, start int) ([]*model.Ranking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRankings", ctx, start)
	ret0, _ := ret[0].([]*model.Ranking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRankings indicates an expected call of ListRankings.
func (mr *MockRankingUseCaseMockRecorder) ListRankings(ctx, start interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRankings", reflect.TypeOf((*MockRankingUseCase)(nil).ListRankings), ctx, start)
}
