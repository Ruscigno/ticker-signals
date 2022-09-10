package insert

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"runtime"

	"github.com/Ruscigno/ticker-signals/internal/api"
	"github.com/Ruscigno/ticker-signals/internal/utils/app"
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
)

func TestInsert(ctx *context.Context, svc *app.Controllers, apiServer *api.TransactionsServiceServer) bool {
	raw, err := ioutil.ReadFile("insert/test-data/account-create-999998-1625468182-1.json")
	if err != nil {
		zap.L().Fatal("unable to read create account file", zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
		return false
	}

	var acc v1.Account
	err = json.Unmarshal(raw, &acc)
	if err != nil {
		zap.L().Fatal("error marshalling create account file", zap.String("file", string(raw)), zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
		return false
	}
	_, err = apiServer.CreateAccount(*ctx, &v1.CreateAccountRequest{
		Account: &acc,
	})
	if err != nil {
		zap.L().Fatal("create account error", zap.String("file", string(raw)), zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
		return false
	}
	aa, err := apiServer.AccSvc.GetByID(acc.AccountId)
	if err != nil {
		zap.L().Fatal("AccSvc.GetByID error", zap.Int64("AccountId", acc.AccountId), zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
		return false
	}
	if aa.AccountID != acc.AccountId {
		zap.L().Fatal("accountId mismatch", zap.Int64("Right AccountId", aa.AccountID), zap.Int64("Left AccountId", acc.AccountId), zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
		return false
	}
	return true
}
