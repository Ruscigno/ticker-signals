// Code generated by MockGen. DO NOT EDIT.
// Source: internal/transaction/deals/repo/interface.go

// package dealsrepo is a generated GoMock package.
package dealsrepo

import (
	reflect "reflect"
	time "time"

	deals "github.com/Ruscigno/ticker-signals/internal/transaction/deals"
	gomock "github.com/golang/mock/gomock"
)

// MockDealsRepository is a mock of DealsRepository interface.
type MockDealsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDealsRepositoryMockRecorder
}

// MockDealsRepositoryMockRecorder is the mock recorder for MockDealsRepository.
type MockDealsRepositoryMockRecorder struct {
	mock *MockDealsRepository
}

// NewMockDealsRepository creates a new mock instance.
func NewMockDealsRepository(ctrl *gomock.Controller) *MockDealsRepository {
	mock := &MockDealsRepository{ctrl: ctrl}
	mock.recorder = &MockDealsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDealsRepository) EXPECT() *MockDealsRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockDealsRepository) Delete(accountID, dealID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", accountID, dealID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDealsRepositoryMockRecorder) Delete(accountID, dealID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDealsRepository)(nil).Delete), accountID, dealID)
}

// GetByID mocks base method.
func (m *MockDealsRepository) GetByID(accountID, dealID int64) (*deals.Deal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", accountID, dealID)
	ret0, _ := ret[0].(*deals.Deal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockDealsRepositoryMockRecorder) GetByID(accountID, dealID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockDealsRepository)(nil).GetByID), accountID, dealID)
}

// GetTickerBeats mocks base method.
func (m *MockDealsRepository) GetTickerBeats(accountID int64, from time.Time) ([]*deals.Deal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTickerBeats", accountID, from)
	ret0, _ := ret[0].([]*deals.Deal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTickerBeats indicates an expected call of GetTickerBeats.
func (mr *MockDealsRepositoryMockRecorder) GetTickerBeats(accountID, from interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTickerBeats", reflect.TypeOf((*MockDealsRepository)(nil).GetTickerBeats), accountID, from)
}

// Insert mocks base method.
func (m *MockDealsRepository) Insert(ac *deals.Deal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ac)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockDealsRepositoryMockRecorder) Insert(ac interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockDealsRepository)(nil).Insert), ac)
}

// Update mocks base method.
func (m *MockDealsRepository) Update(ac *deals.Deal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ac)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockDealsRepositoryMockRecorder) Update(ac interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDealsRepository)(nil).Update), ac)
}

// UpdateOrderIdByPosition mocks base method.
func (m *MockDealsRepository) UpdateOrderIdByPosition(accountID, positionID, orderID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderIdByPosition", accountID, positionID, orderID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderIdByPosition indicates an expected call of UpdateOrderIdByPosition.
func (mr *MockDealsRepositoryMockRecorder) UpdateOrderIdByPosition(accountID, positionID, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderIdByPosition", reflect.TypeOf((*MockDealsRepository)(nil).UpdateOrderIdByPosition), accountID, positionID, orderID)
}