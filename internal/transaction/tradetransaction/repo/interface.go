package tradetransactionrepo

import (
	model "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
)

// TradeTransactionRepository is a CRUD to the database
//
//go:generate mockery -inpkg -name Interface
type TradeTransactionRepository interface {
	InsertTransaction(tt *model.TradeTransaction, timeGMT int64) (int64, error)
	InsertRequest(tt *model.TradeRequest, timeGMT int64) (int64, error)
	InsertResult(tt *model.TradeResult, timeGMT int64) (int64, error)
	GetByIDTransaction(accountID, orderID, creationOrder int64) (*model.TradeTransaction, error)
	GetByIDRequest(accountID, orderID, creationOrder int64) (*model.TradeRequest, error)
	GetByIDResult(accountID, orderID, creationOrder int64) (*model.TradeResult, error)
}
