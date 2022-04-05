package service

import (
	"github.com/kirktriplefive/test"
	"github.com/kirktriplefive/test/pkg/repository"
)

type OrderService struct {
	repo repository.Order
}


func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(order test.Order) (int, error){
	return s.repo.Create(order)
}

func (s *OrderService) CreateOrderWithItem(orderId int, itemId string) (int, error){
	return s.repo.CreateOrderWithItem(orderId, itemId)
}

func (s *OrderService) CreateOrderWithNewItem(orderId int, item test.Item) (int, error){
	return s.repo.CreateOrderWithNewItem(orderId, item)
}

func (s *OrderService) GetById(orderId int) (test.Order, []test.Item, error){
	return s.repo.GetById(orderId)
}