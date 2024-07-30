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

// MockRankingRepository is a mock of RankingRepository interface.
type MockRankingRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRankingRepositoryMockRecorder
}

// MockRankingRepositoryMockRecorder is the mock recorder for MockRankingRepository.
type MockRankingRepositoryMockRecorder struct {
	mock *MockRankingRepository
}

// NewMockRankingRepository creates a new mock instance.
func NewMockRankingRepository(ctrl *gomock.Controller) *MockRankingRepository {
	mock := &MockRankingRepository{ctrl: ctrl}
	mock.recorder = &MockRankingRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRankingRepository) EXPECT() *MockRankingRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRankingRepository) Create(ctx context.Context, key string, ranking *model.Ranking) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, key, ranking)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRankingRepositoryMockRecorder) Create(ctx, key, ranking interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRankingRepository)(nil).Create), ctx, key, ranking)
}

// List mocks base method.
func (m *MockRankingRepository) List(ctx context.Context, key string, start int) ([]*model.Ranking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, key, start)
	ret0, _ := ret[0].([]*model.Ranking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRankingRepositoryMockRecorder) List(ctx, key, start interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRankingRepository)(nil).List), ctx, key, start)
}
