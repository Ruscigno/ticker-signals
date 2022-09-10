package orderssvc

import (
	"context"
	"runtime"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/orders"
	repo "github.com/Ruscigno/ticker-signals/internal/transaction/orders/repo"
	"github.com/blendle/zapdriver"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"

	bb "github.com/Ruscigno/ticker-signals/internal/tickerbeats/beats"
	ss "github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"

	dd "github.com/Ruscigno/ticker-signals/internal/transaction/deals/service"
)

// NewOrdersService creates a service to interact with PostgreSQL
func NewOrdersService(ctx context.Context, repo repo.OrdersRepository, beats bb.TickerBeatsService, deals dd.DealsService) OrdersService {
	return &ordersService{
		ctx:   ctx,
		repo:  repo,
		beats: beats,
		deals: deals,
	}
}

// OrdersService ...
type ordersService struct {
	ctx   context.Context
	repo  repo.OrdersRepository
	beats bb.TickerBeatsService
	deals dd.DealsService
}

// Insert inserts a new order
func (r *ordersService) Insert(ac *model.Order) error {
	err := r.repo.Insert(ac)
	if err != nil {
		zap.L().Error("Order: Insert Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	if ac.Magic.Valid {
		r.beats.ConfirmByExternalID(ac.AccountID, ac.Magic.Int64, ac.Ticket, ac.PositionId, ss.OrderConfirmed)
	}
	//r.deals.UpdateOrderIdByPosition(ac.AccountID, ac.PositionId, ac.OrderID)
	return nil
}

// GetByID return an order by its Id
func (r *ordersService) GetByID(accountID, ticket int64) (*model.Order, error) {
	ac, err := r.repo.GetByID(accountID, ticket)
	if err != nil {
		zap.L().Error("Order: GetByID Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	return ac, nil
}

// Update updates an order
func (r *ordersService) Update(ac *model.Order) error {
	err := r.repo.Update(ac)
	if err != nil {
		zap.L().Error("Order: Update Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

// Delete deletes an order
func (r *ordersService) Delete(accountID, ticket int64) error {
	err := r.repo.Delete(accountID, ticket)
	if err != nil {
		zap.L().Error("Order: Delete Error", zap.Error(err), zap.Int64("accountID", accountID), zap.Int64("order_id", ticket), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}
