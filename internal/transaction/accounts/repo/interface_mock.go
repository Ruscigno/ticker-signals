// Code generated by MockGen. DO NOT EDIT.
// Source: internal/transaction/accounts/repo/interface.go

// package accountsrepo is a generated GoMock package.
package accountsrepo

import (
	reflect "reflect"

	accounts "github.com/Ruscigno/ticker-signals/internal/transaction/accounts"
	gomock "github.com/golang/mock/gomock"
)

// MockAccountRepository is a mock of AccountRepository interface.
type MockAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRepositoryMockRecorder
}

// MockAccountRepositoryMockRecorder is the mock recorder for MockAccountRepository.
type MockAccountRepositoryMockRecorder struct {
	mock *MockAccountRepository
}

// NewMockAccountRepository creates a new mock instance.
func NewMockAccountRepository(ctrl *gomock.Controller) *MockAccountRepository {
	mock := &MockAccountRepository{ctrl: ctrl}
	mock.recorder = &MockAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountRepository) EXPECT() *MockAccountRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockAccountRepository) Delete(accountID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAccountRepositoryMockRecorder) Delete(accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAccountRepository)(nil).Delete), accountID)
}

// GetByID mocks base method.
func (m *MockAccountRepository) GetByID(accountID int64) (*accounts.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", accountID)
	ret0, _ := ret[0].(*accounts.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockAccountRepositoryMockRecorder) GetByID(accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAccountRepository)(nil).GetByID), accountID)
}

// Insert mocks base method.
func (m *MockAccountRepository) Insert(ac *accounts.Account, serverTime int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ac, serverTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockAccountRepositoryMockRecorder) Insert(ac, serverTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockAccountRepository)(nil).Insert), ac, serverTime)
}

// Update mocks base method.
func (m *MockAccountRepository) Update(ac *accounts.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ac)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAccountRepositoryMockRecorder) Update(ac interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAccountRepository)(nil).Update), ac)
}