// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go
//
// Generated by this command:
//
//	mockgen -source=./repository.go -destination=./mocks/mock_repository.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	models "tiger-tracker-api/repository/models"

	gin "github.com/gin-gonic/gin"
	gomock "go.uber.org/mock/gomock"
)

// MockAppRepository is a mock of AppRepository interface.
type MockAppRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAppRepositoryMockRecorder
}

// MockAppRepositoryMockRecorder is the mock recorder for MockAppRepository.
type MockAppRepositoryMockRecorder struct {
	mock *MockAppRepository
}

// NewMockAppRepository creates a new mock instance.
func NewMockAppRepository(ctrl *gomock.Controller) *MockAppRepository {
	mock := &MockAppRepository{ctrl: ctrl}
	mock.recorder = &MockAppRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppRepository) EXPECT() *MockAppRepositoryMockRecorder {
	return m.recorder
}

// GetRecentTigerSightings mocks base method.
func (m *MockAppRepository) GetRecentTigerSightings(ctx *gin.Context, pageNumber, pageSize int) ([]models.TigerDetailWithSightings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecentTigerSightings", ctx, pageNumber, pageSize)
	ret0, _ := ret[0].([]models.TigerDetailWithSightings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecentTigerSightings indicates an expected call of GetRecentTigerSightings.
func (mr *MockAppRepositoryMockRecorder) GetRecentTigerSightings(ctx, pageNumber, pageSize any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecentTigerSightings", reflect.TypeOf((*MockAppRepository)(nil).GetRecentTigerSightings), ctx, pageNumber, pageSize)
}