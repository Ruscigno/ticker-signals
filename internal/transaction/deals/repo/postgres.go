package dealsrepo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Ruscigno/ticker-signals/internal/transaction/deals"
	"github.com/Ruscigno/ticker-signals/internal/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// NewDealsRepo creates a service to interact with PostgreSQL
func NewDealsRepo(ctx context.Context, dbCon *sqlx.DB) DealsRepository {
	return &dealsRepo{
		ctx:   ctx,
		dbCon: dbCon,
	}
}

// dealsrepo ...
type dealsRepo struct {
	ctx   context.Context
	dbCon *sqlx.DB
}

// Insert creates a new deal
func (c *dealsRepo) Insert(ac *deals.Deal) error {
	ac.Created = ac.DealTime
	ac.Updated = time.Now().UTC()
	ac.Deleted.Valid = false
	const Insert string = "INSERT INTO tickerbeats.deals (%s) VALUES(%s) ON CONFLICT ON CONSTRAINT deals_pk DO NOTHING;"
	ra, err := utils.ExecScript(c.ctx, *ac, Insert, c.dbCon, nil)
	if err != nil {
		return err
	}
	if ra > 0 {
		zap.L().Info("Deal inserted", zap.Int64("accountID", ac.AccountID), zap.Int64("ticket", ac.Ticket))
	}
	return err
}

// GetByID gets an deal by Id
func (c *dealsRepo) GetByID(accountID, dealID int64) (*deals.Deal, error) {
	const SelectQuery string = "select %s from tickerbeats.deals where accountID = $1 and dealid = $2 and deleted is null"
	result := deals.Deal{}
	_, fields := utils.StructToSlice(result, utils.DefaultFieldsToIgnore)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	err := c.dbCon.Get(&result, query, accountID, dealID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return &result, nil
		}
		return nil, err
	}
	return &result, nil
}

// Update updates an deal
func (c *dealsRepo) Update(ac *deals.Deal) error {
	return nil
}

// Delete deletes an deal
func (c *dealsRepo) Delete(accountID, dealID int64) error {
	return nil
}

func (c *dealsRepo) GetTickerBeats(accountID int64, from time.Time) ([]*deals.Deal, error) {
	const SelectQuery string = "select %s from tickerbeats.deals where accountID = $1 and dealtime >= $2 and deleted is null"
	result := []*deals.Deal{}
	_, fields := utils.StructToSlice(deals.Deal{}, utils.DefaultFieldsToIgnore)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	err := c.dbCon.Select(&result, query, accountID, from)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return result, nil
		}
		return nil, err
	}
	return result, nil
}

func (c *dealsRepo) UpdateOrderIdByPosition(accountID, positionID, orderID int64) error {
	const UpdateQuery string = "update tickerbeats.deals set orderID = $3 where accountID = $1 and positionID = $2 and orderID is null"
	_, err := c.dbCon.Exec(UpdateQuery, accountID, positionID, orderID)
	return err
}
