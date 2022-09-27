// Code generated by MockGen. DO NOT EDIT.
// Source: internal/transaction/traderules/service/interface.go

// package tradeRulessvc is a generated GoMock package.
package tradeRulessvc

import (
	reflect "reflect"

	tradeRules "github.com/Ruscigno/ticker-signals/internal/transaction/traderules"
	tradetransaction "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
	gomock "github.com/golang/mock/gomock"
)

// MockTradeRulesService is a mock of TradeRulesService interface.
type MockTradeRulesService struct {
	ctrl     *gomock.Controller
	recorder *MockTradeRulesServiceMockRecorder
}

// MockTradeRulesServiceMockRecorder is the mock recorder for MockTradeRulesService.
type MockTradeRulesServiceMockRecorder struct {
	mock *MockTradeRulesService
}

// NewMockTradeRulesService creates a new mock instance.
func NewMockTradeRulesService(ctrl *gomock.Controller) *MockTradeRulesService {
	mock := &MockTradeRulesService{ctrl: ctrl}
	mock.recorder = &MockTradeRulesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTradeRulesService) EXPECT() *MockTradeRulesServiceMockRecorder {
	return m.recorder
}

// ApplyRules mocks base method.
func (m *MockTradeRulesService) ApplyRules(accountID int64, req []*tradetransaction.TradeRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplyRules", accountID, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// ApplyRules indicates an expected call of ApplyRules.
func (mr *MockTradeRulesServiceMockRecorder) ApplyRules(accountID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyRules", reflect.TypeOf((*MockTradeRulesService)(nil).ApplyRules), accountID, req)
}

// GetAll mocks base method.
func (m *MockTradeRulesService) GetAll() ([]*tradeRules.TradeRules, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*tradeRules.TradeRules)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockTradeRulesServiceMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockTradeRulesService)(nil).GetAll))
}

// GetByAccount mocks base method.
func (m *MockTradeRulesService) GetByAccount(accountID int64, ruleType int) ([]*tradeRules.TradeRules, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAccount", accountID, ruleType)
	ret0, _ := ret[0].([]*tradeRules.TradeRules)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByAccount indicates an expected call of GetByAccount.
func (mr *MockTradeRulesServiceMockRecorder) GetByAccount(accountID, ruleType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAccount", reflect.TypeOf((*MockTradeRulesService)(nil).GetByAccount), accountID, ruleType)
}

// GetBySymbol mocks base method.
func (m *MockTradeRulesService) GetBySymbol(accountID int64, symbol string, ruleType int) ([]*tradeRules.TradeRules, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBySymbol", accountID, symbol, ruleType)
	ret0, _ := ret[0].([]*tradeRules.TradeRules)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBySymbol indicates an expected call of GetBySymbol.
func (mr *MockTradeRulesServiceMockRecorder) GetBySymbol(accountID, symbol, ruleType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBySymbol", reflect.TypeOf((*MockTradeRulesService)(nil).GetBySymbol), accountID, symbol, ruleType)
}