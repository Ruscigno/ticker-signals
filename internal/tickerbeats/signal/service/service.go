package signalsvc

import (
	"context"
	"runtime"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	"github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
	repo "github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal/repo"
	tm "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
)

// NewSignalService creates a service to interact with PostgreSQL
func NewSignalService(ctx context.Context, signal repo.SignalRepository) SignalService {
	return &signalService{
		ctx:    ctx,
		signal: signal,
	}
}

// TradeTransactionService ...
type signalService struct {
	ctx    context.Context
	signal repo.SignalRepository
}

// GetSignalByDestination returns a Signal filtered by its `destinationaccountid`
func (r *signalService) GetSignalByDestination(accountID int64) ([]*signal.Signal, error) {
	signal, err := r.signal.GetSignalByDestination(accountID)
	if err != nil {
		zap.L().Error("SignalService: GetSignalByDestination Error",
			zap.Int64("accountID", accountID),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
	return signal, err
}

func (r *signalService) ConfirmByExternalID(destinationAccountID, externalID, tickerBeatsID, positionID int64, status signal.SignalStatusEnum) error {
	err := r.signal.ConfirmByExternalID(destinationAccountID, externalID, tickerBeatsID, positionID, status)
	if err != nil {
		zap.L().Error("SignalService: ConfirmByExternalID Error",
			zap.Int64("destinationAccountID", destinationAccountID),
			zap.Int64("externalID", externalID),
			zap.Int64("tickerBeatsID", tickerBeatsID),
			zap.Int64("positionID", positionID),
			zap.Int64("status", int64(status)),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
	return nil
}

func (r *signalService) UpdateStatus(sourceAccountID, destinationAccountID int64, status signal.SignalStatusEnum, groupID string) error {
	err := r.signal.UpdateStatus(sourceAccountID, destinationAccountID, status, groupID)
	if err != nil {
		zap.L().Error("SignalService: UpdateStatus Error",
			zap.Int64("sourceAccountID", sourceAccountID),
			zap.Int64("destinationAccountID", destinationAccountID),
			zap.Int64("status", int64(status)),
			zap.String("groupID", groupID),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
	return nil
}

func (r *signalService) CreateTickerBeats(sourceAccountID, destinationAccountID int64, groupID string, minToExpire int32) (int64, error) {
	ra, err := r.signal.CreateTickerBeats(sourceAccountID, destinationAccountID, groupID, minToExpire)
	if err != nil {
		zap.L().Error("SignalService: CreateTickerBeats Error",
			zap.Int64("sourceAccountID", sourceAccountID),
			zap.Int64("destinationAccountID", destinationAccountID),
			zap.String("groupID", groupID),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return 0, err
	}
	if ra < 1 {
		return 0, nil
	}
	zap.L().Info("TickerBeats created",
		zap.Int64("RowsAffected", ra),
		zap.Int64("sourceAccountID", sourceAccountID),
		zap.Int64("destinationAccountID", destinationAccountID),
		zap.String("groupID", groupID),
	)
	return ra, nil
}

func (r *signalService) GetTradeRequesByGroupID(destinationAccountID int64, groupID string, entry v1.DealEntry) ([]*tm.TradeRequest, error) {
	trades, err := r.signal.GetTradeRequesByGroupID(destinationAccountID, groupID, entry)
	if err != nil {
		zap.L().Error("SignalService: GetByGroupID Error",
			zap.Int64("destinationAccountID", destinationAccountID),
			zap.String("groupID", groupID),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
	return trades, err
}

func (r *signalService) RemoveDuplicatedSignals(destinationAccountID int64, groupID string) error {
	err := r.signal.RemoveDuplicatedSignals(destinationAccountID, groupID)
	if err != nil {
		zap.L().Error("SignalService: RemoveDuplicatedSignals Error",
			zap.Int64("destinationAccountID", destinationAccountID),
			zap.String("groupID", groupID),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
	return err
}

func (r *signalService) UpdatePositionIdBeforeClose(sourceAccountID, destinationAccountID int64, groupID string) error {
	err := r.signal.UpdatePositionIdBeforeClose(sourceAccountID, destinationAccountID, groupID)
	if err != nil {
		zap.L().Error("SignalService: UpdatePositionIdBeforeClose Error",
			zap.Int64("sourceAccountID", sourceAccountID),
			zap.Int64("destinationAccountID", destinationAccountID),
			zap.String("groupID", groupID),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
	return err
}

func (r *signalService) CloseDeadPositions(sourceAccountID, destinationAccountID int64, groupID string, minToExpire int32) ([]*tm.TradeRequest, error) {
	tr, err := r.signal.CloseDeadPositions(sourceAccountID, destinationAccountID, groupID, minToExpire)
	if err != nil {
		zap.L().Error("SignalService: CloseDeadPositions Error",
			zap.Int64("sourceAccountID", sourceAccountID),
			zap.Int64("destinationAccountID", destinationAccountID),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
	return tr, err

}

func (r *signalService) NeedToCloseAllPositions(accountID int64, groupID string, stop int32) ([]*tm.TradeRequest, error) {
	tr, err := r.signal.NeedToCloseAllPositions(accountID, groupID, stop)
	if err != nil {
		zap.L().Error("SignalService: CloseAllPositions Error",
			zap.Int64("accountID", accountID),
			zap.String("groupID", groupID),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
	return tr, err

}
