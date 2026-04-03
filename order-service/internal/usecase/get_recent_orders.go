package usecase

import (
	"context"
	"fmt"
	"order-service/internal/domain"
	"order-service/internal/repository"
)

type GetRecentOrders struct {
	repo repository.OrderRepository
}

func NewGetRecentOrders(repo repository.OrderRepository) *GetRecentOrders {
	return &GetRecentOrders{}
}

type GetRecentOrderCommand struct {
	Limit int `form:"limit"`
}

func (uc *GetRecentOrders) Execute(ctx context.Context, cmd GetRecentOrderCommand) ([]domain.Order, error) {
	if cmd.Limit <= 1 || cmd.Limit > 100 {
		return nil, fmt.Errorf("limit must be between 1 and 100")

	}

	orders, err := uc.repo.GetRecent(ctx, cmd.Limit)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
