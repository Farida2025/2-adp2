package usecase

import (
	"context"
	"fmt"
	"order-service/internal/domain"
	"order-service/internal/repository"
)

type GetOrder struct {
	repo repository.OrderRepository
}

func NewGetOrder(repo repository.OrderRepository) *GetOrder {
	return &GetOrder{repo: repo}
}

func (uc *GetOrder) Execute(ctx context.Context, id string) (*domain.Order, error) {
	order, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

type CancelOrder struct {
	repo repository.OrderRepository
}

func NewCancelOrder(repo repository.OrderRepository) *CancelOrder {
	return &CancelOrder{repo: repo}
}

func (uc *CancelOrder) Execute(ctx context.Context, id string) error {
	order, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}

	if order.Status != "Pending" {
		return fmt.Errorf("cannot cancel order in status: %s", order.Status)
	}

	return uc.repo.UpdateStatus(ctx, id, "Cancelled")
}
