package accountsInfosvc

import (
	"context"
	"runtime"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/accountsinfo"
	repo "github.com/Ruscigno/ticker-signals/internal/transaction/accountsinfo/repo"
	"github.com/blendle/zapdriver"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
)

// NewAccountsInfoService creates a service to interact with PostgreSQL
func NewAccountsInfoService(ctx context.Context, repo repo.AccountInfoRepository) AccountsInfoService {
	return &accountsInfoService{
		ctx:  ctx,
		repo: repo,
	}
}

// AccountsInfoService ...
type accountsInfoService struct {
	ctx  context.Context
	repo repo.AccountInfoRepository
}

// Insert inserts a new accountInfo
func (r *accountsInfoService) Insert(ac *model.AccountInfo) error {
	err := r.repo.Insert(ac)
	if err != nil {
		zap.L().Error("AccountInfo: Insert Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

// GetByID return an accountInfo by its Id
func (r *accountsInfoService) GetByID(accountID int64, infoID string) (*model.AccountInfo, error) {
	ac, err := r.repo.GetByID(accountID, infoID)
	if err != nil {
		zap.L().Error("AccountInfo: GetByID Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	return ac, nil
}

// Update updates an accountInfo
func (r *accountsInfoService) Update(ac *model.AccountInfo) error {
	err := r.repo.Update(ac)
	if err != nil {
		zap.L().Error("AccountInfo: Update Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

// Delete deletes an accountInfo
func (r *accountsInfoService) Delete(accountID int64, infoID string) error {
	err := r.repo.Delete(accountID, infoID)
	if err != nil {
		zap.L().Error("AccountInfo: Delete Error", zap.Error(err), zap.Int64("accountID", accountID), zap.String("infoID", infoID), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}
