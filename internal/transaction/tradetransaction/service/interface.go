package tradetransactionsvc

import (
	model "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
)

// TradeTransactionService is a CRUD to the database
type TradeTransactionService interface {
	GetByIDRequest(accountID, orderID, creationOrder int64) (*model.TradeRequest, error)
	GetByIDResult(accountID, orderID, creationOrder int64) (*model.TradeResult, error)
	GetByIDTransaction(accountID, orderID, creationOrder int64) (*model.TradeTransaction, error)
	Insert(tt *model.TradeTransaction, tr *model.TradeRequest, rr *model.TradeResult, timeGMT int64) (int64, error)
}
