package accountssvc

import (
	model "github.com/Ruscigno/ticker-signals/internal/transaction/accounts"
)

// AccountsService is a CRUD to the database
//
//go:generate mockery -inpkg -name Interface
type AccountsService interface {
	Insert(ac *model.Account, serverTime int64) error
	GetByID(accountID int64) (*model.Account, error)
	Update(ac *model.Account) error
	Delete(accountID int64) error
}