package orderssvc

import (
	model "github.com/Ruscigno/ticker-signals/internal/transaction/orders"
)

// OrdersService is a CRUD to the database
type OrdersService interface {
	Insert(ac *model.Order) error
	GetByID(accountID, ticket int64) (*model.Order, error)
	Update(ac *model.Order) error
	Delete(accountID, ticket int64) error
}
