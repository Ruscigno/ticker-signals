package accountsInfosvc

import (
	model "github.com/Ruscigno/ticker-signals/internal/transaction/accountsinfo"
)

// AccountsInfoService is a CRUD to the database
//
//go:generate mockery -inpkg -name Interface
type AccountsInfoService interface {
	Insert(ac *model.AccountInfo) error
	GetByID(accountID int64, infoID string) (*model.AccountInfo, error)
	Update(ac *model.AccountInfo) error
	Delete(accountID int64, infoID string) error
}
