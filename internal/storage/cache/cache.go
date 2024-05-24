package cache

import (
	"sync"

	"gihub.com/gmohmad/wb_l0/internal/storage/models/orders"
	uuid "github.com/fossoreslp/go-uuid-v4"
)

type Cache struct {
	lock sync.RWMutex
	Data map[uuid.UUID]orders.OrderItem
}

func NewCache() *Cache {
	return &Cache{
		lock: sync.RWMutex{},
		Data: make(map[uuid.UUID]orders.OrderItem),
	}
}

func (c *Cache) GetOrder(id uuid.UUID) (orders.OrderItem, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	order, ok := c.Data[id]

	return order, ok
}

func (c *Cache) AddOrder(id uuid.UUID, order orders.OrderItem) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.Data[id] = order
}
