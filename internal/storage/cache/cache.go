package cache

import (
	"context"
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

type Storage interface {
	GetOrders(ctx context.Context) ([]orders.Order, error)
}

func (c *Cache) FillUpCache(ctx context.Context, storage Storage) error {
	orders, err := storage.GetOrders(ctx)

	if err != nil {
		return err
	}

	for _, order := range orders {
		c.AddOrder(order.ID, order.OrderItem)
	}

	return nil
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
