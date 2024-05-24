package storage

import (
	"context"
	"fmt"

	"gihub.com/gmohmad/wb_l0/internal/storage/models/orders"
	"gihub.com/gmohmad/wb_l0/internal/storage/postgres"
	uuid "github.com/fossoreslp/go-uuid-v4"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage struct {
	client postgres.Client
}

func NewStorage(client postgres.Client) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) GetOrder(ctx context.Context, id uuid.UUID) (orders.Order, error) {
	query := `SELECT * FROM orders WHERE id = $1`

	var dbId [16]byte
	var ordItem orders.OrderItem

	row := s.client.QueryRow(ctx, query, id)

	if err := row.Scan(&dbId, &ordItem); err != nil {
		return orders.Order{}, err
	}

	id, err := uuid.ParseBytes(dbId[:])

	if err != nil {
		return orders.Order{}, fmt.Errorf("Couldn't parse id returned from db into uuidv4 format: %w", err)
	}

	order := orders.Order{
		ID: id,
		OrderItem: ordItem,
	}

	return order, nil
}

func (s *Storage) SaveOrder(ctx context.Context, order *orders.OrderItem) (*uuid.UUID, error) {

	var dbId [16]byte

	stmt := `INSERT INTO orders (orderItem) VALUES ($1) RETURNING id`

	if err := s.client.QueryRow(ctx, stmt, *order).Scan(&dbId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return nil, fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
		}
		return nil, fmt.Errorf("Couldn't insert order: %w", err)
	}

	id, err := uuid.ParseBytes(dbId[:])

	if err != nil {
		return nil, fmt.Errorf("Couldn't parse id returned from db into uuidv4 format: %w", err)
	}

	return &id, nil
}
