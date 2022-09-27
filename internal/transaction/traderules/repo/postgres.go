package tradeRulesrepo

import (
	"context"
	"fmt"
	"strings"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/traderules"
	"github.com/Ruscigno/ticker-signals/internal/utils"
	"github.com/jmoiron/sqlx"
)

// NewTradeRulesRepo creates a service to interact with PostgreSQL
func NewTradeRulesRepo(ctx context.Context, dbCon *sqlx.DB) TradeRulesRepository {
	return &tradeRulesRepo{
		ctx:   ctx,
		dbCon: dbCon,
	}
}

type tradeRulesRepo struct {
	ctx   context.Context
	dbCon *sqlx.DB
}

// GetByAccount return all rules from an account
func (c *tradeRulesRepo) GetByAccount(accountID int64, ruleType int) ([]*model.TradeRules, error) {
	const SelectQuery string = "select %s from tickerbeats.traderules where accountid=$1 and ruleType=$2 and active=true"

	result := []*model.TradeRules{}
	_, fields := utils.StructToSlice(model.TradeRules{}, nil)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	err := c.dbCon.Select(&result, query, accountID, ruleType)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return result, nil
		}
		return nil, err
	}
	return result, nil
}

// GetBySymbol return all symbol's rules from an account
func (c *tradeRulesRepo) GetBySymbol(accountID int64, symbol string, ruleType int) ([]*model.TradeRules, error) {
	const SelectQuery string = "select %s from tickerbeats.traderules where accountid=$1 and symbol=$2 and ruleType=$3 and active=true"
	result := []*model.TradeRules{}
	_, fields := utils.StructToSlice(model.TradeRules{}, nil)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	err := c.dbCon.Select(&result, query, accountID, symbol, ruleType)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return result, nil
		}
		return nil, err
	}
	return result, nil
}

func (c *tradeRulesRepo) GetAll() ([]*model.TradeRules, error) {
	const SelectQuery string = "select %s from tickerbeats.traderules where active=true"
	result := []*model.TradeRules{}
	_, fields := utils.StructToSlice(model.TradeRules{}, nil)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	err := c.dbCon.Select(&result, query)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return result, nil
		}
		return nil, err
	}
	return result, nil
}
