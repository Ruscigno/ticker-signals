package tradeRulesrepo

import (
	model "github.com/Ruscigno/ticker-signals/internal/transaction/traderules"
)

// TradeRulesRepository is a CRUD to the database
//
//go:generate mockery -inpkg -name Interface
type TradeRulesRepository interface {
	GetByAccount(accountID int64, ruleType int) ([]*model.TradeRules, error)
	GetBySymbol(accountID int64, symbol string, ruleType int) ([]*model.TradeRules, error)
	GetAll() ([]*model.TradeRules, error)
}
