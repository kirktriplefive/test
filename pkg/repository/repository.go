package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kirktriplefive/test"
)

type Authorization interface {
	CreateUser(user test.User) (int, error)
	GetUser(username, password string) (test.User, error)
}

type Item interface {
	Create(item test.Item) (string, error)
	GetAll() ([]test.Item, error)
}

type OrderNew struct {
	order test.Order
	delivery test.Delivery
	payment test.Payment
}

type Order interface {
	Create(order test.Order) (int, error)
	CreateOrderWithItem(orderId int, itemId string) (int, error)
	CreateOrderWithNewItem(orderId int, item test.Item) (int, error)
	GetById(orderId int) (test.Order, []test.Item, error)
}

type Repository struct{
	Authorization
	Item
	Order
}

func NewRepository(db *sqlx.DB) *Repository{
	return &Repository{
		Authorization: NewAuthPostrgres(db),
		Item: NewItemPostrgres(db),
		Order: NewOrderPostgres(db), 

	}
}