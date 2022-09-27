package tradeRulessvc

import (
	model "github.com/Ruscigno/ticker-signals/internal/transaction/traderules"
	tt "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
)

// AccountsInfoService is a CRUD to the database
//
//go:generate mockery -inpkg -name Interface
type TradeRulesService interface {
	ApplyRules(accountID int64, req []*tt.TradeRequest) error
	GetAll() ([]*model.TradeRules, error)
	GetByAccount(accountID int64, ruleType int) ([]*model.TradeRules, error)
	GetBySymbol(accountID int64, symbol string, ruleType int) ([]*model.TradeRules, error)
}
