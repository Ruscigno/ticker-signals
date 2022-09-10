package positionsrepo

import (
	"context"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	model "github.com/Ruscigno/ticker-signals/internal/transaction/positions"
	"github.com/Ruscigno/ticker-signals/internal/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// NewPositionsRepo creates a service to interact with PostgreSQL
func NewPositionsRepo(ctx context.Context, dbCon *sqlx.DB) PositionsRepository {
	return &positionRepo{
		ctx:   ctx,
		dbCon: dbCon,
	}
}

// positionRepo ...
type positionRepo struct {
	ctx   context.Context
	dbCon *sqlx.DB
}

// Insert creates a new account
func (c *positionRepo) Insert(ac *model.Position) (int64, error) {
	ac.Created = time.Now().UTC()
	ac.Updated = ac.Created
	if ac.PositionTime.Valid {
		ac.Created = ac.PositionTime.Time
	}
	ac.Deleted.Valid = false
	Insert := "INSERT INTO tickerbeats.positions (%s) VALUES(%s) ON CONFLICT ON CONSTRAINT positions_pk DO NOTHING;"
	ra, err := utils.ExecScript(c.ctx, *ac, Insert, c.dbCon, nil)
	if err != nil {
		return 0, err
	}
	if ra == 0 {
		return 0, nil
	}
	zap.L().Info("Position inserted",
		zap.Int64("accountID", ac.AccountID),
		zap.Int64("ticket", ac.Ticket))
	return ra, nil
}

// GetByID gets an account by Id
func (c *positionRepo) GetByID(accountID, ticket int64) (*model.Position, error) {
	const SelectQuery string = "select %s from tickerbeats.positions where accountID = $1 and ticket = $2"
	result := model.Position{}
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
func (c *positionRepo) Update(ac *model.Position) error {
	const UpdateQuery string = `
        update tickerbeats.positions set 
            pricecurrent = %f, 
            volume = %f, 
            stoploss = %f, 
            takeprofit = %f, 
            commission = %f, 
            swap = %f, 
            profit = %f, 
            updated = current_timestamp 
            where accountID = %d and ticket = %d
        `
	sql := fmt.Sprintf(UpdateQuery, ac.PriceCurrent, ac.Volume, ac.StopLoss, ac.TakeProfit, ac.Commission, ac.Swap, ac.Profit, ac.AccountID, ac.Ticket)
	_, err := utils.ExecScript(c.ctx, nil, sql, c.dbCon, utils.DefaultFieldsToIgnore)
	if err != nil {
		return err
	}
	// zap.L().Debug("Position updated",
	// 	zap.Int64("accountID", ac.AccountID),
	// 	zap.Int64("ticket", ac.Ticket))
	return nil
}

func (c *positionRepo) CloseOne(accountID, ticket int64, closed time.Time, commission, swap, profit float64) error {
	const DeleteQuery string = "update tickerbeats.positions set deleted=$1,updated=$1,commission=$4,swap=$5,profit=$6 where accountid=$2 and ticket=$3 and deleted is null"
	rs, err := c.dbCon.Exec(DeleteQuery, closed, accountID, ticket, commission, swap, profit)
	if err != nil {
		return err
	}
	ra, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if ra > 0 {
		zap.L().Info("Position closed",
			zap.Int64("accountID", accountID),
			zap.Int64("ticket", ticket),
			zap.Time("closed", closed),
		)
	}
	return nil
}

func (c *positionRepo) CloseAll(accountID int64, closed time.Time, commission, swap, profit float64) error {
	const DeleteQuery string = "update tickerbeats.positions set deleted = $1, updated = $1 where accountid = $2 and positiontime < $3"
	rs, err := c.dbCon.Exec(DeleteQuery, closed, accountID)
	if err != nil {
		return err
	}
	ra, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if ra > 0 {
		zap.L().Info("All positions closed",
			zap.Int64("accountID", accountID),
			zap.Time("closed", closed),
		)
	}
	return nil
}

func (c *positionRepo) CloseIfNotIn(accountID int64, tickets []int64, closed, maxDate time.Time) error {
	if len(tickets) == 0 {
		return nil
	}
	tkts := make([]string, len(tickets))
	for i, value := range tickets {
		tkts[i] = fmt.Sprintf("%d", value)
	}
	tktList := strings.Join(tkts[:], ",")
	const DeleteQuery string = "update tickerbeats.positions set deleted = $1, updated = $1 where accountid = $2 and ticket not in (%s) and deleted is null and positiontime <= $3"
	rs, err := c.dbCon.Exec(fmt.Sprintf(DeleteQuery, tktList), closed, accountID, maxDate)
	if err != nil {
		return err
	}
	ra, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if ra > 0 {
		zap.L().Info("Positions closed",
			zap.Int64("accountID", accountID),
			zap.String("tickets", tktList),
			zap.Time("closed", closed),
		)
	}
	return nil
}

func (c *positionRepo) GetTickerBeats(sourceAccountID, destinationAccountID int64, from time.Time) ([]*model.Position, error) {
	const SelectQuery string = `
    SELECT o.* 
    FROM tickerbeats.positions o,
         tickerbeats.signals s
    WHERE o.accountid = s.sourceaccountid
        and s.sourceaccountid = $1
        and s.destinationaccountid = $2
        and s.active = true
        and o.positiontime  >= $3::TIMESTAMP - (INTERVAL '1 min' * s.minutestoexpire)
        and o.ticket not in (select sr.sourcebeatsid 
                             from tickerbeats.signalsresult sr 
                             where sr.signaltype = $4
                                   and sr.sourceaccountid = s.sourceaccountid 
                                   and sr.destinationaccountid = s.destinationaccountid
                                   and sr.confirmationtime is not null)
`

	result := []*model.Position{}
	err := c.dbCon.Select(&result, SelectQuery, sourceAccountID, destinationAccountID, from, v1.SignalType_name[int32(v1.SignalType_SIGINAL_TYPE_POSITION)])
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return result, nil
		}
		return nil, err
	}
	return result, nil
}

func (c *positionRepo) ConfirmByExternalID(destinationAccountID, externalID int64, beatTime time.Time) error {
	const UpdateQuery string = "update tickerbeats.signalsresult set confirmationtime = $1 where destinationAccountID = $2 and externalID = $3"
	_, err := c.dbCon.Exec(UpdateQuery, beatTime, destinationAccountID, externalID)
	if err != nil {
		return err
	}
	zap.L().Info("Position confirmed",
		zap.Int64("destinationAccountID", destinationAccountID),
		zap.Int64("externalID", externalID),
		zap.Time("beatTime", beatTime),
	)
	return nil
}

func (c *positionRepo) GetActivePositions(accountID int64) ([]*model.Position, error) {
	_, fields := utils.StructToSlice(model.Position{}, utils.DefaultFieldsToIgnore)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql := psql.Select().
		Columns(fields...).
		From(model.TableName).
		Where(sq.And{
			sq.Eq{"accountid": accountID},
			sq.NotEq{"deleted": nil},
		}).
		OrderBy("positionid")

	query, _, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	result := []*model.Position{}
	err = c.dbCon.Select(&result, query)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return result, nil
		}
		return nil, err
	}
	return result, nil
}
