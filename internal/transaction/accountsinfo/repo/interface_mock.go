// Code generated by MockGen. DO NOT EDIT.
// Source: internal/transaction/accountsinfo/repo/interface.go

// package accountsInforepo is a generated GoMock package.
package accountsInforepo

import (
	reflect "reflect"
	time "time"

	accountsInfo "github.com/Ruscigno/ticker-signals/internal/transaction/accountsinfo"
	gomock "github.com/golang/mock/gomock"
)

// MockAccountInfoRepository is a mock of AccountInfoRepository interface.
type MockAccountInfoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountInfoRepositoryMockRecorder
}

// MockAccountInfoRepositoryMockRecorder is the mock recorder for MockAccountInfoRepository.
type MockAccountInfoRepositoryMockRecorder struct {
	mock *MockAccountInfoRepository
}

// NewMockAccountInfoRepository creates a new mock instance.
func NewMockAccountInfoRepository(ctrl *gomock.Controller) *MockAccountInfoRepository {
	mock := &MockAccountInfoRepository{ctrl: ctrl}
	mock.recorder = &MockAccountInfoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountInfoRepository) EXPECT() *MockAccountInfoRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockAccountInfoRepository) Delete(accountID int64, infoID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", accountID, infoID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAccountInfoRepositoryMockRecorder) Delete(accountID, infoID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAccountInfoRepository)(nil).Delete), accountID, infoID)
}

// GetByID mocks base method.
func (m *MockAccountInfoRepository) GetByID(accountID int64, infoID string) (*accountsInfo.AccountInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", accountID, infoID)
	ret0, _ := ret[0].(*accountsInfo.AccountInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockAccountInfoRepositoryMockRecorder) GetByID(accountID, infoID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAccountInfoRepository)(nil).GetByID), accountID, infoID)
}

// GetByTimeGMT mocks base method.
func (m *MockAccountInfoRepository) GetByTimeGMT(accountID int64, timeGMT time.Time) (*accountsInfo.AccountInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTimeGMT", accountID, timeGMT)
	ret0, _ := ret[0].(*accountsInfo.AccountInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTimeGMT indicates an expected call of GetByTimeGMT.
func (mr *MockAccountInfoRepositoryMockRecorder) GetByTimeGMT(accountID, timeGMT interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTimeGMT", reflect.TypeOf((*MockAccountInfoRepository)(nil).GetByTimeGMT), accountID, timeGMT)
}

// Insert mocks base method.
func (m *MockAccountInfoRepository) Insert(ac *accountsInfo.AccountInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ac)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockAccountInfoRepositoryMockRecorder) Insert(ac interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockAccountInfoRepository)(nil).Insert), ac)
}

// Update mocks base method.
func (m *MockAccountInfoRepository) Update(ac *accountsInfo.AccountInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ac)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAccountInfoRepositoryMockRecorder) Update(ac interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAccountInfoRepository)(nil).Update), ac)
}