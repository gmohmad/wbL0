package storage

import (
	"context"
	"fmt"

	"gihub.com/gmohmad/wb_l0/internal/storage/models/orders"
	"gihub.com/gmohmad/wb_l0/internal/storage/postgres"
	uuid "github.com/fossoreslp/go-uuid-v4"
	"github.com/jackc/pgx/v5"
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

	row := s.client.QueryRow(ctx, query, id)

	return scanOrder(row)
}

func (s *Storage) GetOrders(ctx context.Context) ([]orders.Order, error) {
	query := `SELECT * FROM orders`

	rows, err := s.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Error querying database: %w", err)
	}
	defer rows.Close()

	var orderSlice []orders.Order
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orderSlice = append(orderSlice, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Row iteration error: %w", err)
	}

	return orderSlice, nil
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

func scanOrder(row pgx.Row) (orders.Order, error) {
	var dbId [16]byte
	var ordItem orders.OrderItem

	if err := row.Scan(&dbId, &ordItem); err != nil {
		return orders.Order{}, err
	}

	id, err := uuid.ParseBytes(dbId[:])
	if err != nil {
		return orders.Order{}, fmt.Errorf("Couldn't parse id returned from db into uuidv4 format: %w", err)
	}

	return orders.Order{
		ID:        id,
		OrderItem: ordItem,
	}, nil
}
