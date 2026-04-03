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

func (r *OrderRepo) GetRecent(ctx context.Context, limit int) ([]domain.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, customer_id, item_name, amount, status, created_at 
         FROM orders 
         ORDER BY created_at DESC 
         LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		err := rows.Scan(&o.ID, &o.CustomerID, &o.ItemName)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)

	}
	return orders, nil
}
