package cache

import (
	"sync"

	uuid "github.com/fossoreslp/go-uuid-v4"
)

type Cache struct {
	lock sync.RWMutex
	Data map[uuid.UUID]interface{}
}

func NewCache() *Cache {
	return &Cache{
		lock: sync.RWMutex{},
		Data: make(map[uuid.UUID]interface{}),
	}
}

func (c *Cache) GetOrder(id uuid.UUID) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.Unlock()

	order, ok := c.Data[id]


	return order, ok
}

func (c *Cache) AddOrder(id uuid.UUID, order interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.Data[id] = order
}
