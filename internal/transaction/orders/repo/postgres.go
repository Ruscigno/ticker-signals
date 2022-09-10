package ordersrepo

import (
	"context"
	"fmt"
	"strings"
	"time"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/orders"
	"github.com/Ruscigno/ticker-signals/internal/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// NewOrdersRepo creates a service to interact with PostgreSQL
func NewOrdersRepo(ctx context.Context, dbCon *sqlx.DB) OrdersRepository {
	return &orderRepo{
		ctx:   ctx,
		dbCon: dbCon,
	}
}

// orderRepo ...
type orderRepo struct {
	ctx   context.Context
	dbCon *sqlx.DB
}

// Insert creates a new account
func (c *orderRepo) Insert(ac *model.Order) error {
	ac.Created = ac.TimeSetup.Time
	ac.Updated = time.Now().UTC()
	ac.Deleted.Valid = false
	const Insert string = "INSERT INTO tickerbeats.orders (%s) VALUES(%s) ON CONFLICT ON CONSTRAINT orders_pk DO NOTHING;"
	ra, err := utils.ExecScript(c.ctx, *ac, Insert, c.dbCon, nil)
	if err != nil {
		return err
	}
	if ra > 0 {
		zap.L().Info("Order inserted", zap.Int64("accountID", ac.AccountID), zap.Int64("ticket", ac.Ticket))
	}
	return err
}

// GetByID gets an account by Id
func (c *orderRepo) GetByID(accountID, ticket int64) (*model.Order, error) {
	const SelectQuery string = "select %s from tickerbeats.orders where accountID = $1 and ticket = $2 and deleted is null"
	result := model.Order{}
	_, fields := utils.StructToSlice(result, utils.DefaultFieldsToIgnore)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	err := c.dbCon.Get(&result, query, accountID, ticket)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return &result, nil
		}
		return nil, err
	}
	return &result, nil
}

// Update updates an account
func (c *orderRepo) Update(ac *model.Order) error {
	return nil
}

// Delete deletes an account
func (c *orderRepo) Delete(accountID, ticket int64) error {
	return nil
}
