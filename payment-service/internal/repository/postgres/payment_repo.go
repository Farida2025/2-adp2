package postgres

import (
	"context"
	"database/sql"
	"payment-service/internal/domain"
)

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Save(ctx context.Context, p domain.Payment) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO payments (id, order_id, transaction_id, amount, status) VALUES ($1, $2, $3, $4, $5)",
		p.ID, p.OrderID, p.TransactionID, p.Amount, p.Status)
	return err
}

func (r *PaymentRepo) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.QueryRowContext(ctx,
		"SELECT id, order_id, transaction_id, amount, status, created_at FROM payments WHERE order_id = $1",
		orderID).Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}
