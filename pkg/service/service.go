package service

import (
	"github.com/gin-gonic/gin"
	"github.com/kirktriplefive/test"
	"github.com/kirktriplefive/test/pkg/cache"
	"github.com/kirktriplefive/test/pkg/repository"
)

type Order interface {
	Create(order test.Order) (int, error)
	CreateOrderWithItem(orderId int, itemId string) (int, error)
	CreateOrderWithNewItem(orderId int, item test.Item) (int, error)
	GetById(orderId int) (test.Order, []test.Item, error)
	CreateNewOrder(order test.Order, items []test.Item) (string, error)
	Close()
	CacheFromPQ() error
	GetOrderById(c *gin.Context) 
}

type Service struct {
	Order
}


func NewService(repos *repository.Repository, cache *cache.Cache) *Service{
	return &Service{
		Order: &OrderService{*repos, *cache},

	}

}