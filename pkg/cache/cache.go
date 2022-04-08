package cache

import (
	"errors"
	"sync"

	"github.com/kirktriplefive/test"
	"github.com/sirupsen/logrus"
)

type Cache struct {
	orders map[string]*test.OrderResponseCache
	mutex *sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]*test.OrderResponseCache),
		mutex: new(sync.RWMutex),
	}
}

func (cache *Cache) AddOrder(order test.OrderResponseCache) {
	cache.orders[order.Order_uid] = &order
}

func (cache *Cache) AddToCache(order test.OrderResponseCache) {
	cache.mutex.Lock()
	cache.orders[order.Order_uid] = &order
	cache.mutex.Unlock()
}

func (cache *Cache) GetOrder(order_uid string) (*test.OrderResponseCache, error) {
	cache.mutex.RLock()
	order, ok := cache.orders[order_uid]
	if !ok {
		logrus.Printf("Order not found")
		return nil, errors.New("Order not found")
	}
	cache.mutex.RUnlock()
	return order, nil
}