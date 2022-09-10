package ordersrepo

import (
	model "github.com/Ruscigno/ticker-signals/internal/transaction/orders"
)

// OrdersRepository is a CRUD to the database
//
//go:generate mockery -inpkg -name Interface
type OrdersRepository interface {
	Insert(ac *model.Order) error
	GetByID(accountID, ticket int64) (*model.Order, error)
	Update(ac *model.Order) error
	Delete(accountID, ticket int64) error
}
