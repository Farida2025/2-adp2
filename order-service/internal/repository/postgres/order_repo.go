package postgres

import (
	"context"
	"database/sql"
	"order-service/internal/domain"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Save(ctx context.Context, o domain.Order) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO orders (id, customer_id, item_name, amount, status) VALUES ($1,$2,$3,$4,$5)",
		o.ID, o.CustomerID, o.ItemName, o.Amount, o.Status)
	return err
}

func (r *OrderRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var o domain.Order
	err := r.db.QueryRowContext(ctx,
		"SELECT id, customer_id, item_name, amount, status, created_at FROM orders WHERE id = $1",
		id).Scan(&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &o, err
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE id = $2", status, id)
	return err
}
