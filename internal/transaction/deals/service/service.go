package dealssvc

import (
	"context"
	"runtime"
	"time"

	"github.com/Ruscigno/ticker-signals/internal/transaction/deals"
	repo "github.com/Ruscigno/ticker-signals/internal/transaction/deals/repo"
	ps "github.com/Ruscigno/ticker-signals/internal/transaction/positions/service"
	"github.com/blendle/zapdriver"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"

	bb "github.com/Ruscigno/ticker-signals/internal/tickerbeats/beats"
	ss "github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
)

// NewDealsService creates a service to interact with PostgreSQL
func NewDealsService(ctx context.Context, repo repo.DealsRepository, pos ps.PositionsService, beats bb.TickerBeatsService) DealsService {
	return &dealsService{
		ctx:   ctx,
		repo:  repo,
		pos:   pos,
		beats: beats,
	}
}

// DealsService ...
type dealsService struct {
	ctx   context.Context
	repo  repo.DealsRepository
	pos   ps.PositionsService
	beats bb.TickerBeatsService
}

// Insert inserts a new deal
func (r *dealsService) Insert(ac *deals.Deal) error {
	err := r.repo.Insert(ac)
	if err != nil {
		zap.L().Error("Deal Service: Insert Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	var commission float64 = 0
	var swap float64 = 0
	var profit float64 = 0
	if ac.Commission.Valid {
		commission = ac.Commission.Float64
	}
	if ac.Swap.Valid {
		commission = ac.Swap.Float64
	}
	if ac.Profit.Valid {
		commission = ac.Profit.Float64
	}
	if ac.Entry == v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_OUT)] {
		r.pos.CloseOne(ac.AccountID, ac.PositionId, ac.DealTime, commission, swap, profit)
	}
	if ac.Magic.Valid {
		r.beats.ConfirmByExternalID(ac.AccountID, ac.Magic.Int64, ac.Ticket, ac.PositionId, ss.DealConfirmed)
	}
	return nil
}

// GetByID return an deal by its Id
func (r *dealsService) GetByID(accountID, dealID int64) (*deals.Deal, error) {
	ac, err := r.repo.GetByID(accountID, dealID)
	if err != nil {
		zap.L().Error("Deal Service: GetByID Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	return ac, nil
}

// Update updates an deal
func (r *dealsService) Update(ac *deals.Deal) error {
	err := r.repo.Update(ac)
	if err != nil {
		zap.L().Error("Deal Service: Update Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

// Delete deletes an deal
func (r *dealsService) Delete(accountID, dealID int64) error {
	err := r.repo.Delete(accountID, dealID)
	if err != nil {
		zap.L().Error("Deal Service: Delete Error", zap.Error(err), zap.Int64("accountID", accountID), zap.Int64("deal_id", dealID), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

// GetTickerBeats return an ticker-beats deals from an account ID
func (r *dealsService) GetTickerBeats(accountID int64, from time.Time) ([]*deals.Deal, error) {
	deals, err := r.repo.GetTickerBeats(accountID, from)
	if err != nil {
		zap.L().Error("Deal Service: GetTickerBeats Error",
			zap.Int64("accountID", accountID),
			zap.Time("from", from),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	return deals, nil
}

func (r *dealsService) UpdateOrderIdByPosition(accountID, positionID, orderID int64) error {
	err := r.repo.UpdateOrderIdByPosition(accountID, positionID, orderID)
	if err != nil {
		zap.L().Error("Deal Service: UpdateOrderIdByPosition Error",
			zap.Int64("accountID", accountID),
			zap.Int64("positionID", positionID),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}
