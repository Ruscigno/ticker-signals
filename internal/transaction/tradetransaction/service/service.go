package tradetransactionsvc

import (
	"context"
	"runtime"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
	repo "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction/repo"
	"github.com/blendle/zapdriver"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"

	bb "github.com/Ruscigno/ticker-signals/internal/tickerbeats/beats"
	ss "github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
)

// NewTradeTransactionService creates a service to interact with PostgreSQL
func NewTradeTransactionService(ctx context.Context, repo repo.TradeTransactionRepository, beats bb.TickerBeatsService) TradeTransactionService {
	return &tradeTransactionService{
		ctx:   ctx,
		repo:  repo,
		beats: beats,
	}
}

// TradeTransactionService ...
type tradeTransactionService struct {
	ctx   context.Context
	repo  repo.TradeTransactionRepository
	beats bb.TickerBeatsService
}

// Insert inserts a new tradetransaction
func (r *tradeTransactionService) Insert(tt *model.TradeTransaction, tr *model.TradeRequest, rr *model.TradeResult, timeGMT int64) (int64, error) {
	var ratt, rarr, ratr int64

	//TODO: run all of them in ONE transaction
	ratt, err := r.repo.InsertTransaction(tt, timeGMT)
	if err == nil {
		ratr, err = r.repo.InsertRequest(tr, timeGMT)
	}
	if err == nil {
		rarr, err = r.repo.InsertResult(rr, timeGMT)
	}
	if err != nil {
		zap.L().Error("TradeTransaction: Insert Error",
			zap.String("TradeTransaction.model", spew.Sdump(tt)),
			zap.String("TradeRequest.model", spew.Sdump(tr)),
			zap.String("TradeResult.model", spew.Sdump(rr)),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return 0, err
	}
	if tr.Magic.Valid {
		var positionID int64 = 0
		if tr.PositionID.Valid {
			positionID = tr.PositionID.Int64
		}
		r.beats.ConfirmByExternalID(tr.AccountID, tr.Magic.Int64, tr.OrderID, positionID, ss.TransactionConfirmed)
	}
	return ratt + rarr + ratr, nil
}

func (r *tradeTransactionService) GetByIDTransaction(accountID, orderID, creationOrder int64) (*model.TradeTransaction, error) {
	ac, err := r.repo.GetByIDTransaction(accountID, orderID, creationOrder)
	if err != nil {
		zap.L().Error("GetByIDTransaction Error",
			zap.Int64("accountID", accountID),
			zap.Int64("order_id", orderID),
			zap.Int64("creation_order", creationOrder),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	return ac, nil
}

func (r *tradeTransactionService) GetByIDRequest(accountID, orderID, creationOrder int64) (*model.TradeRequest, error) {
	ac, err := r.repo.GetByIDRequest(accountID, orderID, creationOrder)
	if err != nil {
		zap.L().Error("GetByIDRequest Error",
			zap.Int64("accountID", accountID),
			zap.Int64("order_id", orderID),
			zap.Int64("creation_order", creationOrder),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	return ac, nil
}

func (r *tradeTransactionService) GetByIDResult(accountID, orderID, creationOrder int64) (*model.TradeResult, error) {
	ac, err := r.repo.GetByIDResult(accountID, orderID, creationOrder)
	if err != nil {
		zap.L().Error("GetByIDResult Error",
			zap.Int64("accountID", accountID),
			zap.Int64("order_id", orderID),
			zap.Int64("creation_order", creationOrder),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	return ac, nil
}
