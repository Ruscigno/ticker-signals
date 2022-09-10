package accountsInforepo

import (
	"time"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/accountsinfo"
)

// AccountInfoRepository is a CRUD to the database
//
//go:generate mockery -inpkg -name Interface
type AccountInfoRepository interface {
	Insert(ac *model.AccountInfo) error
	GetByID(accountID int64, infoID string) (*model.AccountInfo, error)
	GetByTimeGMT(accountID int64, timeGMT time.Time) (*model.AccountInfo, error)
	Update(ac *model.AccountInfo) error
	Delete(accountID int64, infoID string) error
}
