package dealssvc

import (
	"time"

	"github.com/Ruscigno/ticker-signals/internal/transaction/deals"
)

// DealsService is a CRUD to the database
type DealsService interface {
	Insert(ac *deals.Deal) error
	GetByID(accountID, dealID int64) (*deals.Deal, error)
	Update(ac *deals.Deal) error
	Delete(accountID, dealID int64) error
	GetTickerBeats(accountID int64, from time.Time) ([]*deals.Deal, error)
	UpdateOrderIdByPosition(accountID, positionID, orderID int64) error
}
