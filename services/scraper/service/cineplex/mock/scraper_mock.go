// Code generated by MockGen. DO NOT EDIT.
// Source: ./scraper.go

// Package mock_cineplex is a generated GoMock package.
package mock_cineplex

import (
	reflect "reflect"
	dto "scraper/dto"
	model "scraper/storage/model"

	gomock "github.com/golang/mock/gomock"
)

// MockFilmStorageRepository is a mock of FilmStorageRepository interface.
type MockFilmStorageRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFilmStorageRepositoryMockRecorder
}

// MockFilmStorageRepositoryMockRecorder is the mock recorder for MockFilmStorageRepository.
type MockFilmStorageRepositoryMockRecorder struct {
	mock *MockFilmStorageRepository
}

// NewMockFilmStorageRepository creates a new mock instance.
func NewMockFilmStorageRepository(ctrl *gomock.Controller) *MockFilmStorageRepository {
	mock := &MockFilmStorageRepository{ctrl: ctrl}
	mock.recorder = &MockFilmStorageRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFilmStorageRepository) EXPECT() *MockFilmStorageRepositoryMockRecorder {
	return m.recorder
}

// Insert mocks base method.
func (m *MockFilmStorageRepository) Insert(film model.Film) (model.Film, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", film)
	ret0, _ := ret[0].(model.Film)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockFilmStorageRepositoryMockRecorder) Insert(film interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockFilmStorageRepository)(nil).Insert), film)
}

// IsExists mocks base method.
func (m *MockFilmStorageRepository) IsExists(film model.Film) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExists", film)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsExists indicates an expected call of IsExists.
func (mr *MockFilmStorageRepositoryMockRecorder) IsExists(film interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExists", reflect.TypeOf((*MockFilmStorageRepository)(nil).IsExists), film)
}

// MockSoup is a mock of Soup interface.
type MockSoup struct {
	ctrl     *gomock.Controller
	recorder *MockSoupMockRecorder
}

// MockSoupMockRecorder is the mock recorder for MockSoup.
type MockSoupMockRecorder struct {
	mock *MockSoup
}

// NewMockSoup creates a new mock instance.
func NewMockSoup(ctrl *gomock.Controller) *MockSoup {
	mock := &MockSoup{ctrl: ctrl}
	mock.recorder = &MockSoupMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSoup) EXPECT() *MockSoupMockRecorder {
	return m.recorder
}

// GetMovies mocks base method.
func (m *MockSoup) GetMovies() ([]dto.Film, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovies")
	ret0, _ := ret[0].([]dto.Film)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovies indicates an expected call of GetMovies.
func (mr *MockSoupMockRecorder) GetMovies() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovies", reflect.TypeOf((*MockSoup)(nil).GetMovies))
}
