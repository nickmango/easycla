// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
//

// Code generated by MockGen. DO NOT EDIT.
// Source: events/repository.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/communitybridge/easycla/cla-backend-go/gen/v1/models"
	events "github.com/communitybridge/easycla/cla-backend-go/gen/v1/restapi/operations/events"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddDataToEvent mocks base method.
func (m *MockRepository) AddDataToEvent(eventID, parentProjectSFID, projectSFID, projectSFName, companySFID, projectID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDataToEvent", eventID, parentProjectSFID, projectSFID, projectSFName, companySFID, projectID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDataToEvent indicates an expected call of AddDataToEvent.
func (mr *MockRepositoryMockRecorder) AddDataToEvent(eventID, parentProjectSFID, projectSFID, projectSFName, companySFID, projectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDataToEvent", reflect.TypeOf((*MockRepository)(nil).AddDataToEvent), eventID, parentProjectSFID, projectSFID, projectSFName, companySFID, projectID)
}

// CreateEvent mocks base method.
func (m *MockRepository) CreateEvent(event *models.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEvent", event)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEvent indicates an expected call of CreateEvent.
func (mr *MockRepositoryMockRecorder) CreateEvent(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEvent", reflect.TypeOf((*MockRepository)(nil).CreateEvent), event)
}

// GetClaGroupEvents mocks base method.
func (m *MockRepository) GetClaGroupEvents(claGroupID string, nextKey *string, paramPageSize *int64, all bool, searchTerm *string) (*models.EventList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClaGroupEvents", claGroupID, nextKey, paramPageSize, all, searchTerm)
	ret0, _ := ret[0].(*models.EventList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClaGroupEvents indicates an expected call of GetClaGroupEvents.
func (mr *MockRepositoryMockRecorder) GetClaGroupEvents(claGroupID, nextKey, paramPageSize, all, searchTerm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClaGroupEvents", reflect.TypeOf((*MockRepository)(nil).GetClaGroupEvents), claGroupID, nextKey, paramPageSize, all, searchTerm)
}

// GetCompanyClaGroupEvents mocks base method.
func (m *MockRepository) GetCompanyClaGroupEvents(claGroupID, companySFID string, nextKey *string, paramPageSize *int64, searchTerm *string, all bool) (*models.EventList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyClaGroupEvents", claGroupID, companySFID, nextKey, paramPageSize, searchTerm, all)
	ret0, _ := ret[0].(*models.EventList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyClaGroupEvents indicates an expected call of GetCompanyClaGroupEvents.
func (mr *MockRepositoryMockRecorder) GetCompanyClaGroupEvents(claGroupID, companySFID, nextKey, paramPageSize, searchTerm, all interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyClaGroupEvents", reflect.TypeOf((*MockRepository)(nil).GetCompanyClaGroupEvents), claGroupID, companySFID, nextKey, paramPageSize, searchTerm, all)
}

// GetCompanyEvents mocks base method.
func (m *MockRepository) GetCompanyEvents(companyID, eventType string, nextKey *string, paramPageSize *int64, all bool) (*models.EventList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyEvents", companyID, eventType, nextKey, paramPageSize, all)
	ret0, _ := ret[0].(*models.EventList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyEvents indicates an expected call of GetCompanyEvents.
func (mr *MockRepositoryMockRecorder) GetCompanyEvents(companyID, eventType, nextKey, paramPageSize, all interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyEvents", reflect.TypeOf((*MockRepository)(nil).GetCompanyEvents), companyID, eventType, nextKey, paramPageSize, all)
}

// GetCompanyFoundationEvents mocks base method.
func (m *MockRepository) GetCompanyFoundationEvents(companySFID, companyID, foundationSFID string, nextKey *string, paramPageSize *int64, searchTerm *string, all bool) (*models.EventList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyFoundationEvents", companySFID, companyID, foundationSFID, nextKey, paramPageSize, searchTerm, all)
	ret0, _ := ret[0].(*models.EventList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyFoundationEvents indicates an expected call of GetCompanyFoundationEvents.
func (mr *MockRepositoryMockRecorder) GetCompanyFoundationEvents(companySFID, companyID, foundationSFID, nextKey, paramPageSize, searchTerm, all interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyFoundationEvents", reflect.TypeOf((*MockRepository)(nil).GetCompanyFoundationEvents), companySFID, companyID, foundationSFID, nextKey, paramPageSize, searchTerm, all)
}

// GetFoundationEvents mocks base method.
func (m *MockRepository) GetFoundationEvents(foundationSFID string, nextKey *string, paramPageSize *int64, all bool, searchTerm *string) (*models.EventList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFoundationEvents", foundationSFID, nextKey, paramPageSize, all, searchTerm)
	ret0, _ := ret[0].(*models.EventList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFoundationEvents indicates an expected call of GetFoundationEvents.
func (mr *MockRepositoryMockRecorder) GetFoundationEvents(foundationSFID, nextKey, paramPageSize, all, searchTerm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFoundationEvents", reflect.TypeOf((*MockRepository)(nil).GetFoundationEvents), foundationSFID, nextKey, paramPageSize, all, searchTerm)
}

// GetRecentEvents mocks base method.
func (m *MockRepository) GetRecentEvents(pageSize int64) (*models.EventList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecentEvents", pageSize)
	ret0, _ := ret[0].(*models.EventList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecentEvents indicates an expected call of GetRecentEvents.
func (mr *MockRepositoryMockRecorder) GetRecentEvents(pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecentEvents", reflect.TypeOf((*MockRepository)(nil).GetRecentEvents), pageSize)
}

// SearchEvents mocks base method.
func (m *MockRepository) SearchEvents(params *events.SearchEventsParams, pageSize int64) (*models.EventList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchEvents", params, pageSize)
	ret0, _ := ret[0].(*models.EventList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchEvents indicates an expected call of SearchEvents.
func (mr *MockRepositoryMockRecorder) SearchEvents(params, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchEvents", reflect.TypeOf((*MockRepository)(nil).SearchEvents), params, pageSize)
}