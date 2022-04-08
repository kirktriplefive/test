package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirktriplefive/test"
	"github.com/kirktriplefive/test/pkg/cache"
	"github.com/kirktriplefive/test/pkg/repository"
)

type OrderService struct {
	repo repository.Order
	cache cache.Cache
}


func NewOrderService(repo repository.Repository, cache cache.Cache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
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

func (s *OrderService) CreateNewOrder(order test.Order, items []test.Item) (string, error){
	if _,err := s.repo.CreateNewOrder(order, items); err != nil {
		return " ", err
	} else {
		orderResponse := test.OrderResponseCache{
			Order_uid:         order.Order_uid,
			TrackNumber:       order.TrackNumber,
			Entry:             order.Entry,
			Del:               order.Delivery,
			Paym:              order.Payment,
			Items:             items,
			Locale:            order.Locale,
			InternalSignature: order.InternalSignature,
			CustomerId:        order.CustomerId,
			DeliveryService:   order.DeliveryService,
			ShardKey:          order.ShardKey,
			SmId:              order.SmId,
			DateCreated:       order.DateCreated,
			OofShard:          order.OofShard,
		}
		s.cache.AddOrder(orderResponse)
	} 
	return " ", nil
}

func (s *OrderService) CacheFromPQ() error {
	orders, err := s.repo.GetOrdersForCache()
	if err != nil {
		return err
	}
	for _, order := range orders {
		s.cache.AddOrder(order)
	}
	return nil
}

func (s *OrderService) Close() {
	s.repo.Close()
}

func (s *OrderService) GetOrderById(c *gin.Context) {
	order_uid := fmt.Sprintf("%v", c.Param("id"))
	order, err := s.cache.GetOrder(order_uid)
	if err != nil {
		c.HTML(
			http.StatusOK,
			"order.html",
			gin.H{
				"Order_uid" : "Не найден",
			},
		)
		
	} else {
		c.HTML(
			http.StatusOK,
			"order.html",
			gin.H{
				"Order_uid" : order.Order_uid,
				"TrackNumber" :   order.TrackNumber,
				"Entry" : order.Entry  ,
				"Delivery":   order.Del,
				"Payment" :   order.Paym,
				"payload" :   order.Items,
				"Locale" :   order.Locale,
				"InternalSignature" :   order.InternalSignature,
				"CustomerId" :   order.CustomerId,
				"DeliveryService":   order.DeliveryService,
				"ShardKey" :   order.ShardKey,
				"SmId":   order.SmId,
				"DateCreated" :   order.DateCreated,
				"OofShard" :   order.OofShard,
			},
		)
	}



	
}