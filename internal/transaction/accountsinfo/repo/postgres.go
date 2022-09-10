package accountsInforepo

import (
	"context"
	"fmt"
	"strings"
	"time"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/accountsinfo"
	"github.com/Ruscigno/ticker-signals/internal/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// NewAccountInfoRepo creates a service to interact with PostgreSQL
func NewAccountInfoRepo(ctx context.Context, dbCon *sqlx.DB) AccountInfoRepository {
	return &accountsInfoRepo{
		ctx:   ctx,
		dbCon: dbCon,
		// nextTime: make(map[int64]time.Time),
	}
}

type accountsInfoRepo struct {
	ctx   context.Context
	dbCon *sqlx.DB
	// nextTime map[int64]time.Time
}

const (
	Insert     = "INSERT INTO tickerbeats.accountsInfo (%s) VALUES(%s) ON CONFLICT ON CONSTRAINT accountsinfo_pk DO NOTHING;"
	SelectByID = "select %s from tickerbeats.accountsinfo where accountid = $1 and infoid = $2"
)

// Insert creates a new accountInfo
func (c *accountsInfoRepo) Insert(ac *model.AccountInfo) error {
	// next, ok := c.nextTime[ac.AccountID]
	// if ok && time.Now().UTC().Before(next) {
	// 	return nil
	// }
	timeGMT := time.Now().UTC()
	if ac.TimeGMT.Valid {
		timeGMT = ac.TimeGMT.Time.UTC()
	}
	account, err := c.GetByTimeGMT(ac.AccountID, timeGMT)
	if err != nil {
		return err
	}
	if account.AccountID == ac.AccountID {
		return nil
	}
	var seq int64
	err = c.dbCon.Get(&seq, "select nextval('tickerbeats.accountsinfo_seq') as value")
	if err != nil {
		return err
	}
	ac.InfoID = seq
	_, err = utils.ExecScript(c.ctx, *ac, Insert, c.dbCon, nil)
	// if ra > 0 {
	// 	zap.L().Info("Account Info inserted", zap.Int64("accountID", ac.AccountID), zap.Int64("infoID", ac.InfoID))
	// }
	//TODO: fix it
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return nil
	}
	// c.nextTime[ac.AccountID] = time.Now().UTC().Add(time.Minute)
	return err
}

// GetByID gets an accountInfo by Id
func (c *accountsInfoRepo) GetByID(accountID int64, infoID string) (*model.AccountInfo, error) {
	result := model.AccountInfo{}
	_, fields := utils.StructToSlice(result, utils.DefaultFieldsToIgnore)
	query := fmt.Sprintf(SelectByID, strings.Join(fields[:], ","))
	err := c.dbCon.Get(&result, query, accountID, infoID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return &result, nil
		}
		zap.L().Error("AccountInfo: GetByID", zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return &result, nil
}

// Update updates an accountInfo
func (c *accountsInfoRepo) Update(ac *model.AccountInfo) error {
	// ac.Updated = time.Now().UTC()
	return nil
}

// Delete deletes an accountInfo
func (c *accountsInfoRepo) Delete(accountID int64, infoID string) error {
	// ac.Deleted = time.Now().UTC()
	return nil
}

func (c *accountsInfoRepo) GetByTimeGMT(accountID int64, timeGMT time.Time) (*model.AccountInfo, error) {
	const SelectQuery string = "select %s from tickerbeats.accountsInfo where accountID = $1 and timeGMT = $2"
	result := []model.AccountInfo{}
	_, fields := utils.StructToSlice(model.AccountInfo{}, utils.DefaultFieldsToIgnore)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	// zap.L().Info("aaa", zap.String("query", query), zap.String("fields", strings.Join(fields[:], ",")))
	err := c.dbCon.Select(&result, query, accountID, timeGMT)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return &model.AccountInfo{}, nil
		}
		return nil, err
	}
	if len(result) != 1 {
		return &model.AccountInfo{}, nil
	}
	return &result[0], nil
}
