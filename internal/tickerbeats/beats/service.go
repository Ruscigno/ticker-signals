package tickerbeats

import (
	"context"
	"time"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	"github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
	ss "github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
	sSvc "github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal/service"
	tm "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
	"github.com/google/uuid"
	"go.uber.org/zap"

	tr "github.com/Ruscigno/ticker-signals/internal/transaction/traderules/service"
)

// NewTickerBeatsService creates a service to interact with PostgreSQL
func NewTickerBeatsService(ctx context.Context, signal sSvc.SignalService, rules tr.TradeRulesService) TickerBeatsService {
	return &tickerBeatsService{
		ctx:    ctx,
		signal: signal,
		rules:  rules,
	}
}

// TradeTransactionService ...
type tickerBeatsService struct {
	ctx    context.Context
	signal sSvc.SignalService
	rules  tr.TradeRulesService
}

func (r *tickerBeatsService) ConfirmByExternalID(destinationAccountID, externalid, tickerBeatsID, positionID int64, status signal.SignalStatusEnum) error {
	err := r.signal.ConfirmByExternalID(destinationAccountID, externalid, tickerBeatsID, positionID, status)
	if err != nil {
		zap.L().Error("Confirm order error",
			zap.Int64("destinationAccountID", destinationAccountID),
			zap.Int64("externalid", externalid),
			zap.Int64("tickerBeatsID", tickerBeatsID),
			zap.Int64("positionID", positionID),
			zap.Int64("status", int64(status)),
			zap.Error(err),
		)
	}
	return err
}

func (r *tickerBeatsService) GetTickerBeats(accountID int64, from time.Time, sType v1.SignalType, lastOrderID int64) ([]*TickerBeats, error) {
	var result []*TickerBeats
	signals, err := r.signal.GetSignalByDestination(accountID)
	if err != nil {
		return nil, err
	}
	for _, signal := range signals {
		if sType != v1.SignalType_SIGINAL_TYPE_GET_ALL && sType != v1.SignalType_SIGINAL_TYPE_TRADE_REQUEST {
			continue
		}

		var tradesIN, trades []*tm.TradeRequest
		groupID := uuid.New().String()
		ra, err := r.signal.CreateTickerBeats(signal.SourceAccountID, accountID, groupID, signal.MinutesToExpire)
		if err != nil {
			return nil, err
		}
		// Position enter / Deals IN
		if ra > 0 {
			tradesIN, err = r.getBeatsFromPositions(accountID, groupID, signal)
			if err != nil {
				return nil, err
			}
			trades = append(trades, tradesIN...)
		}
		//Close dead positions
		deads, err := r.signal.CloseDeadPositions(signal.SourceAccountID, accountID, groupID, signal.MinutesToExpire)
		if err != nil {
			return nil, err
		}
		trades = append(trades, deads...)

		//Adding trades to the final result
		if len(trades) == 0 {
			continue
		}
		trades, err = r.applyRules(accountID, trades)
		if err != nil {
			return nil, err
		}
		if len(trades) == 0 {
			continue
		}
		calculateOrderBoost(signal, trades)
		ticker := &TickerBeats{
			Signal:         signal,
			ExpirationTime: signal.MinutesToExpire,
			GroupID:        groupID,
			TradeRequests:  trades,
			Valid:          true,
		}
		result = append(result, ticker)
	}
	return result, nil
}

func (r *tickerBeatsService) applyRules(accountID int64, trades []*tm.TradeRequest) ([]*tm.TradeRequest, error) {
	err := r.rules.ApplyRules(accountID, trades)
	if err != nil {
		return nil, err
	}
	var result []*tm.TradeRequest
	for _, trade := range trades {
		if trade.Status != tm.Invalid {
			result = append(result, trade)
		}
	}
	return result, nil
}

func calculateOrderBoost(signal *signal.Signal, trades []*tm.TradeRequest) {
	if !signal.OrderBoost.Valid {
		return
	}
	if !signal.OrderBoostType.Valid {
		return
	}
	otype := ss.OrderBoostTypeEnum(signal.OrderBoostType.Int32)
	oValue := signal.OrderBoost.Float64
	for _, t := range trades {
		if otype == ss.Multiplier {
			t.Volume = t.Volume * oValue
		} else if otype == ss.OrderLimit {
			t.Volume = oValue
		}
	}
}

func (r *tickerBeatsService) BeatsSent(beats []*TickerBeats) error {
	for _, bb := range beats {
		err := r.signal.UpdateStatus(bb.Signal.SourceAccountID, bb.Signal.DestinationAccountID, ss.Sent, bb.GroupID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *tickerBeatsService) getBeatsFromPositions(accountID int64, groupID string, signal *ss.Signal) ([]*tm.TradeRequest, error) {
	err := r.signal.UpdatePositionIdBeforeClose(signal.SourceAccountID, accountID, groupID)
	if err != nil {
		return nil, err
	}
	err = r.signal.RemoveDuplicatedSignals(accountID, groupID)
	if err != nil {
		return nil, err
	}
	tradesIN, err := r.signal.GetTradeRequesByGroupID(accountID, groupID, v1.DealEntry_DEAL_ENTRY_IN)
	return tradesIN, err
}
