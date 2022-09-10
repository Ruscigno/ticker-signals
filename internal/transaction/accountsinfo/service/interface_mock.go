// Code generated by MockGen. DO NOT EDIT.
// Source: internal/transaction/accountsinfo/service/interface.go

// package accountsInfosvc is a generated GoMock package.
package accountsInfosvc

import (
	reflect "reflect"

	accountsInfo "github.com/Ruscigno/ticker-signals/internal/transaction/accountsinfo"
	gomock "github.com/golang/mock/gomock"
)

// MockAccountsInfoService is a mock of AccountsInfoService interface.
type MockAccountsInfoService struct {
	ctrl     *gomock.Controller
	recorder *MockAccountsInfoServiceMockRecorder
}

// MockAccountsInfoServiceMockRecorder is the mock recorder for MockAccountsInfoService.
type MockAccountsInfoServiceMockRecorder struct {
	mock *MockAccountsInfoService
}

// NewMockAccountsInfoService creates a new mock instance.
func NewMockAccountsInfoService(ctrl *gomock.Controller) *MockAccountsInfoService {
	mock := &MockAccountsInfoService{ctrl: ctrl}
	mock.recorder = &MockAccountsInfoServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountsInfoService) EXPECT() *MockAccountsInfoServiceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockAccountsInfoService) Delete(accountID int64, infoID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", accountID, infoID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAccountsInfoServiceMockRecorder) Delete(accountID, infoID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAccountsInfoService)(nil).Delete), accountID, infoID)
}

// GetByID mocks base method.
func (m *MockAccountsInfoService) GetByID(accountID int64, infoID string) (*accountsInfo.AccountInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", accountID, infoID)
	ret0, _ := ret[0].(*accountsInfo.AccountInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockAccountsInfoServiceMockRecorder) GetByID(accountID, infoID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAccountsInfoService)(nil).GetByID), accountID, infoID)
}

// Insert mocks base method.
func (m *MockAccountsInfoService) Insert(ac *accountsInfo.AccountInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ac)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockAccountsInfoServiceMockRecorder) Insert(ac interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockAccountsInfoService)(nil).Insert), ac)
}

// Update mocks base method.
func (m *MockAccountsInfoService) Update(ac *accountsInfo.AccountInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ac)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAccountsInfoServiceMockRecorder) Update(ac interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAccountsInfoService)(nil).Update), ac)
}