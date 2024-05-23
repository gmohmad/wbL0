package storage

import (
	"context"

	"gihub.com/gmohmad/wb_l0/internal/storage/postgres"
	uuid "github.com/fossoreslp/go-uuid-v4"
)

type Storage struct {
	client postgres.Client
}

func (s *Storage) GetOrder(ctx context.Context, id uuid.UUID) (interface{}, error) 

func (s *Storage) GetOrders(ctx context.Context) ([]interface{}, error)

func (s *Storage) SaveOrder(ctx context.Context, id uuid.UUID, order interface{}) error
