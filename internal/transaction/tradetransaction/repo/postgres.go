package tradetransactionrepo

import (
	"context"
	"fmt"
	"strings"
	"time"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
	"github.com/Ruscigno/ticker-signals/internal/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var defaultFieldsToIgnore = []string{"status"}

// NewTradeTransactionRepo creates a service to interact with PostgreSQL
func NewTradeTransactionRepo(ctx context.Context, dbCon *sqlx.DB) TradeTransactionRepository {
	defaultFieldsToIgnore = append(defaultFieldsToIgnore, utils.DefaultFieldsToIgnore...)
	return &tradetransactionRepo{
		ctx:   ctx,
		dbCon: dbCon,
	}
}

// tradetransactionRepo ...
type tradetransactionRepo struct {
	ctx   context.Context
	dbCon *sqlx.DB
}

func (c *tradetransactionRepo) InsertTransaction(tt *model.TradeTransaction, timeGMT int64) (int64, error) {
	if tt.AccountID <= 0 || tt.OrderID <= 0 {
		return 0, nil
	}
	ttm, err := c.GetByIDTransaction(tt.AccountID, tt.OrderID, tt.CreationOrder)
	if err != nil {
		return 0, err
	} else if ttm.AccountID == tt.AccountID && ttm.OrderID == tt.OrderID && ttm.CreationOrder == tt.CreationOrder {
		zap.L().Info("TradeTransaction already exist",
			zap.Int64("accountID", tt.AccountID),
			zap.Int64("order_id", tt.OrderID),
			zap.Int64("creation_order", tt.CreationOrder))
		return 0, nil
	}
	tt.Updated = time.Now().UTC()
	tt.Deleted.Valid = false
	tt.Created = time.Unix(timeGMT, 0)
	if tt.Created.Before(time.Unix(0, utils.MT5_MIN_DATE)) {
		tt.Created = tt.Updated
	}
	validateFields(tt)
	const InsertTransaction string = "INSERT INTO tickerbeats.tradetransactions (%s,internalid) VALUES(%s,nextval('tickerbeats.tradetransactions_seq')) ON CONFLICT ON CONSTRAINT TradeTransactions_pk DO NOTHING;"
	ra, err := utils.ExecScript(c.ctx, *tt, InsertTransaction, c.dbCon, []string{"internalid"})
	if err != nil {
		return 0, err
	}
	if ra > 0 {
		zap.L().Info("TradeTransaction inserted",
			zap.Int64("accountID", tt.AccountID),
			zap.Int64("order_id", tt.OrderID),
			zap.Int64("creation_order", tt.CreationOrder))
	}
	return ra, nil
}

func validateFields(tt *model.TradeTransaction) {
	if strings.Contains(tt.TradeType, "TRADE_TRANSACTION_ORDER") {
		tt.DealType = ""
		return
	}
	if strings.Contains(tt.TradeType, "TRADE_TRANSACTION_DEAL") {
		tt.OrderType = ""
		tt.OrderState = ""
		return
	}
	if strings.Contains(tt.TradeType, "TRADE_TRANSACTION_HISTORY") {
		tt.DealType = ""
		return
	}
}

func (c *tradetransactionRepo) InsertRequest(tt *model.TradeRequest, timeGMT int64) (int64, error) {
	if tt.AccountID <= 0 || tt.OrderID <= 0 {
		return 0, nil
	}
	ttm, err := c.GetByIDRequest(tt.AccountID, tt.OrderID, tt.CreationOrder)
	if err != nil {
		return 0, err
	} else if ttm.AccountID == tt.AccountID && ttm.OrderID == tt.OrderID && ttm.CreationOrder == tt.CreationOrder {
		zap.L().Info("TradeRequest already exist",
			zap.Int64("accountID", tt.AccountID),
			zap.Int64("order_id", tt.OrderID),
			zap.Int64("creation_order", tt.CreationOrder))
		return 0, nil
	}
	tt.Updated = time.Now().UTC()
	tt.Deleted.Valid = false
	tt.Created = time.Unix(timeGMT, 0)
	if tt.Created.Before(time.Unix(0, utils.MT5_MIN_DATE)) {
		tt.Created = tt.Updated
	}
	const InsertRequest string = "INSERT INTO tickerbeats.traderequests (%s) VALUES(%s) ON CONFLICT ON CONSTRAINT TradeRequests_pk DO NOTHING;"
	ra, err := utils.ExecScript(c.ctx, *tt, InsertRequest, c.dbCon, []string{"status"})
	if err != nil {
		return 0, err
	}
	if ra > 0 {
		zap.L().Info("TradeRequest inserted",
			zap.Int64("accountID", tt.AccountID),
			zap.Int64("order_id", tt.OrderID),
			zap.Int64("creation_order", tt.CreationOrder))
	}
	return ra, nil
}

func (c *tradetransactionRepo) InsertResult(tt *model.TradeResult, timeGMT int64) (int64, error) {
	if tt.AccountID <= 0 || tt.OrderID <= 0 {
		return 0, nil
	}
	tt.Updated = time.Now().UTC()
	tt.Deleted.Valid = false
	tt.Created = time.Unix(timeGMT, 0)
	if tt.Created.Before(time.Unix(0, utils.MT5_MIN_DATE)) {
		tt.Created = tt.Updated
	}
	const InsertResult string = "INSERT INTO tickerbeats.traderesults (%s) VALUES(%s) ON CONFLICT ON CONSTRAINT TradeResults_pk DO NOTHING;"
	ra, err := utils.ExecScript(c.ctx, *tt, InsertResult, c.dbCon, nil)
	if err != nil {
		return 0, err
	}
	if ra > 0 {
		zap.L().Info("TradeResult inserted",
			zap.Int64("accountID", tt.AccountID),
			zap.Int64("order_id", tt.OrderID),
			zap.Int64("creation_order", tt.CreationOrder))
	}
	return 0, nil
}

func (c *tradetransactionRepo) GetByIDTransaction(accountID, orderID, creationOrder int64) (*model.TradeTransaction, error) {
	const SelectQuery string = "select %s from tickerbeats.tradetransactions where AccountID = $1 and OrderID = $2 and CreationOrder = $3 and deleted is null"
	result := model.TradeTransaction{}
	_, fields := utils.StructToSlice(result, defaultFieldsToIgnore)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	err := c.dbCon.Get(&result, query, accountID, orderID, creationOrder)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return &result, nil
		}
		return nil, err
	}
	return &result, nil
}

func (c *tradetransactionRepo) GetByIDRequest(accountID, orderID, creationOrder int64) (*model.TradeRequest, error) {
	const SelectQuery string = "select %s from tickerbeats.traderequests where AccountID = $1 and OrderID = $2 and CreationOrder = $3 and deleted is null"
	result := model.TradeRequest{}
	_, fields := utils.StructToSlice(result, defaultFieldsToIgnore)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	err := c.dbCon.Get(&result, query, accountID, orderID, creationOrder)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return &result, nil
		}
		return nil, err
	}
	return &result, nil
}

func (c *tradetransactionRepo) GetByIDResult(accountID, orderID, creationOrder int64) (*model.TradeResult, error) {
	const SelectQuery string = "select %s from tickerbeats.traderesults where AccountID = $1 and OrderID = $2 and CreationOrder = $3 and deleted is null"
	result := model.TradeResult{}
	_, fields := utils.StructToSlice(result, defaultFieldsToIgnore)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	err := c.dbCon.Get(&result, query, accountID, orderID, creationOrder)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return &result, nil
		}
		return nil, err
	}
	return &result, nil
}
