package accountssvc

import (
	"context"
	"runtime"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/accounts"
	repo "github.com/Ruscigno/ticker-signals/internal/transaction/accounts/repo"
	"github.com/blendle/zapdriver"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
)

// NewAccountsService creates a service to interact with PostgreSQL
func NewAccountsService(ctx context.Context, repo repo.AccountRepository) AccountsService {
	return &accountsService{
		ctx:  ctx,
		repo: repo,
	}
}

// AccountsService ...
type accountsService struct {
	ctx  context.Context
	repo repo.AccountRepository
}

// Insert inserts a new account
func (r *accountsService) Insert(ac *model.Account, serverTime int64) error {
	err := r.repo.Insert(ac, serverTime)
	if err != nil {
		zap.L().Error("Account: Insert Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

// GetByID return an account by its Id
func (r *accountsService) GetByID(accountID int64) (*model.Account, error) {
	ac, err := r.repo.GetByID(accountID)
	if err != nil {
		zap.L().Error("Account: GetByID Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	return ac, nil
}

// Update updates an account
func (r *accountsService) Update(ac *model.Account) error {
	err := r.repo.Update(ac)
	if err != nil {
		zap.L().Error("Account: Update Error", zap.Error(err), zap.String("model", spew.Sdump(ac)), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}

// Delete deletes an account
func (r *accountsService) Delete(accountID int64) error {
	err := r.repo.Delete(accountID)
	if err != nil {
		zap.L().Error("Account: Delete Error", zap.Int64("accountID", accountID), zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
		return err
	}
	return nil
}
