package service

import (
	"github.com/kirktriplefive/test"
	"github.com/kirktriplefive/test/pkg/repository"
)

type Authorization interface {
	CreateUser(user test.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Item interface {
	Create(item test.Item) (string, error)
	GetAll() ([]test.Item, error)
}

type Order interface {
	Create(order test.Order) (int, error)
	CreateOrderWithItem(orderId int, itemId string) (int, error)
	CreateOrderWithNewItem(orderId int, item test.Item) (int, error)
	GetById(orderId int) (test.Order, []test.Item, error)
}

type Service struct {
	Authorization
	Item
	Order
}


func NewService(repos *repository.Repository) *Service{
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Item: NewItemService(repos.Item),
		Order: NewOrderService(repos.Order),

	}

}