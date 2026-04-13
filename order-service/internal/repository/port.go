package repository

import (
	"context"
	"order-service/internal/domain"
)

type OrderRepository interface {
	Save(ctx context.Context, o domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id, status string) error
}
