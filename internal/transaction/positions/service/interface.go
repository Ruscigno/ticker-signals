package positionssvc

import (
	"time"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/positions"
)

// PositionsService is a CRUD to the database
type PositionsService interface {
	Insert(ac *model.Position) (int64, error)
	InsertMulti(accountId int64, ac []*model.Position) (bool, error)
	GetByID(accountID, ticket int64) (*model.Position, error)
	Update(ac *model.Position) error
	CloseIfNotIn(accountID int64, tickets []int64, closed, maxDate time.Time) error
	CloseOne(accountID, ticket int64, closed time.Time, commission, swap, profit float64) error
	CloseAll(accountID int64, closed time.Time, commission, swap, profit float64) error
	GetTickerBeats(sourceAccountID, destinationAccountID int64, from time.Time) ([]*model.Position, error)
	ConfirmByExternalID(destinationAccountID, externalID int64, beatTime time.Time) error
	GetActivePositions(accountID int64) ([]*model.Position, error)
}
