package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kirktriplefive/test"
)


type Order interface {
	Create(order test.Order) (int, error)
	CreateOrderWithItem(orderId int, itemId string) (int, error)
	CreateOrderWithNewItem(orderId int, item test.Item) (int, error)
	GetById(orderId int) (test.Order, []test.Item, error)
	CreateNewOrder(order test.Order, items []test.Item) (string, error)
	Close() error
	GetOrdersForCache() ([]test.OrderResponseCache, error)
}

type Repository struct{
	Order
}

func NewRepository(db *sqlx.DB) *Repository{
	return &Repository{
		Order: NewOrderPostgres(db), 

	}
}