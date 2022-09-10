package positionssvc

import (
	"context"
	"runtime"
	"time"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/positions"
	repo "github.com/Ruscigno/ticker-signals/internal/transaction/positions/repo"
	"github.com/blendle/zapdriver"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"

	bb "github.com/Ruscigno/ticker-signals/internal/tickerbeats/beats"
	ss "github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
)

// NewPositionsService creates a service to interact with PostgreSQL
func NewPositionsService(ctx context.Context, repo repo.PositionsRepository, beats bb.TickerBeatsService) PositionsService {
	return &positionsService{
		ctx:   ctx,
		repo:  repo,
		beats: beats,
	}
}

// PositionsService ...
type positionsService struct {
	ctx   context.Context
	repo  repo.PositionsRepository
	beats bb.TickerBeatsService
}

// Insert inserts a new position
func (r *positionsService) Insert(ac *model.Position) (int64, error) {
	ra, err := r.repo.Insert(ac)
	if err != nil {
		zap.L().Error("Position: Insert Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return 0, err
	}
	if ac.Magic.Valid {
		r.beats.ConfirmByExternalID(ac.AccountID, ac.Magic.Int64, 0, ac.Ticket, ss.PositionConfirmed)
	}
	return ra, nil
}

// Insert inserts a new position
func (r *positionsService) InsertMulti(accountId int64, ac []*model.Position) (bool, error) {
	tickets := make([]int64, len(ac))
	var maxDate time.Time = time.Unix(0, 0)
	inserted := false
	for i, pos := range ac {
		ra, err := r.repo.Insert(pos)
		if err == nil && ra == 0 {
			err = r.repo.Update(pos)
		}
		if err != nil {
			zap.L().Error("Position: Insert Error", zap.Error(err), zap.String("model", spew.Sdump(pos)), zapdriver.SourceLocation(runtime.Caller(0)))
			return false, err
		}
		inserted = ra > 0
		if pos.PositionTime.Valid && pos.PositionTime.Time.After(maxDate) {
			maxDate = pos.PositionTime.Time
		}
		tickets[i] = pos.Ticket
		if pos.Magic.Valid {
			r.beats.ConfirmByExternalID(pos.AccountID, pos.Magic.Int64, 0, pos.Ticket, ss.PositionConfirmed)
		}
	}
	r.CloseIfNotIn(accountId, tickets, time.Now().UTC(), maxDate)
	return inserted, nil
}

// GetByID return an position by its Id
func (r *positionsService) GetByID(accountID, ticket int64) (*model.Position, error) {
	ac, err := r.repo.GetByID(accountID, ticket)
	if err != nil {
		zap.L().Error("Position: GetByID Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	return ac, nil
}

// Update updates an position
func (r *positionsService) Update(ac *model.Position) error {
	err := r.repo.Update(ac)
	if err != nil {
		zap.L().Error("Position: Update Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

func (r *positionsService) CloseOne(accountID, ticket int64, closed time.Time, commission, swap, profit float64) error {
	err := r.repo.CloseOne(accountID, ticket, closed, commission, swap, profit)
	if err != nil {
		zap.L().Error("Position: CloseOne Error", zap.Error(err), zap.Int64("accountID", accountID), zap.Int64("ticket", ticket), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

func (r *positionsService) CloseAll(accountID int64, closed time.Time, commission, swap, profit float64) error {
	err := r.repo.CloseAll(accountID, closed, commission, swap, profit)
	if err != nil {
		zap.L().Error("Position: CloseOne Error", zap.Error(err), zap.Int64("accountID", accountID), zap.Time("closed", closed), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

func (r *positionsService) CloseIfNotIn(accountID int64, tickets []int64, closed, maxDate time.Time) error {
	err := r.repo.CloseIfNotIn(accountID, tickets, closed, maxDate)
	if err != nil {
		zap.L().Error("Position: CloseIfNotIn Error", zap.Error(err), zap.Int64("accountID", accountID), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

// GetTickerBeats return an ticker-beats positions from an account ID
func (r *positionsService) GetTickerBeats(sourceAccountID, destinationAccountID int64, from time.Time) ([]*model.Position, error) {
	pos, err := r.repo.GetTickerBeats(sourceAccountID, destinationAccountID, from)
	if err != nil {
		zap.L().Error("Position Service: GetTickerBeats Error",
			zap.Int64("sourceAccountID", sourceAccountID),
			zap.Int64("destinationAccountID", destinationAccountID),
			zap.Time("from", from),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	return pos, nil
}

func (r *positionsService) ConfirmByExternalID(destinationAccountID, externalID int64, beatTime time.Time) error {
	err := r.repo.ConfirmByExternalID(destinationAccountID, externalID, beatTime)
	if err != nil {
		zap.L().Error("Position Service: GetTickerBeats Error",
			zap.Int64("destinationAccountID", destinationAccountID),
			zap.Int64("externalID", externalID),
			zap.Time("beatTime", beatTime),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

func (r *positionsService) GetActivePositions(accountID int64) ([]*model.Position, error) {
	pos, err := r.repo.GetActivePositions(accountID)
	if err != nil {
		zap.L().Error("Position Service: GetActivePositions Error",
			zap.Int64("accountID", accountID),
			zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	return pos, nil
}
