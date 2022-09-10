package dealsrepo

import (
	"time"

	"github.com/Ruscigno/ticker-signals/internal/transaction/deals"
)

// DealsRepository is a CRUD to the database
//
//go:generate mockery -inpkg -name Interface
type DealsRepository interface {
	Insert(ac *deals.Deal) error
	GetByID(accountID, dealID int64) (*deals.Deal, error)
	Update(ac *deals.Deal) error
	Delete(accountID, dealID int64) error
	GetTickerBeats(accountID int64, from time.Time) ([]*deals.Deal, error)
	UpdateOrderIdByPosition(accountID, positionID, orderID int64) error
}
