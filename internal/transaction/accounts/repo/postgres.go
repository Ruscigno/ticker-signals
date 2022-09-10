package accountsrepo

import (
	"context"
	"fmt"
	"strings"
	"time"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/accounts"
	"github.com/Ruscigno/ticker-signals/internal/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// NewAccountsRepo creates a service to interact with PostgreSQL
func NewAccountsRepo(ctx context.Context, dbCon *sqlx.DB) AccountRepository {
	return &accountsRepo{
		ctx:   ctx,
		dbCon: dbCon,
	}
}

// AccountsRepo ...
type accountsRepo struct {
	ctx   context.Context
	dbCon *sqlx.DB
}

// Insert creates a new account
func (c *accountsRepo) Insert(ac *model.Account, serverTime int64) error {
	timeGMT := time.Unix(serverTime, 0)
	ac.Created = time.Now().UTC()
	ac.Updated = ac.Created
	ac.Deleted.Valid = false
	if serverTime > utils.MT5_MIN_DATE {
		ac.Created = timeGMT
	}
	const Insert string = "INSERT INTO tickerbeats.accounts (%s) VALUES(%s) ON CONFLICT ON CONSTRAINT accounts_pk DO NOTHING;"
	ra, err := utils.ExecScript(c.ctx, *ac, Insert, c.dbCon, nil)
	if ra > 0 {
		zap.L().Info("Account inserted", zap.Int64("accountID", ac.AccountID))
	}
	return err
}

// GetByID gets an account by Id
func (c *accountsRepo) GetByID(accountID int64) (*model.Account, error) {
	const SelectQuery string = "select %s from tickerbeats.accounts where accountID = $1"
	result := model.Account{}
	_, fields := utils.StructToSlice(result, utils.DefaultFieldsToIgnore)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	err := c.dbCon.Get(&result, query, accountID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return &result, nil
		}
		return nil, err
	}
	return &result, nil
}

// Update updates an account
func (c *accountsRepo) Update(ac *model.Account) error {
	ac.Updated = time.Now().UTC()
	return nil
}

// Delete deletes an account
func (c *accountsRepo) Delete(accountID int64) error {
	// ac.Deleted = time.Now().UTC()
	return nil
}
