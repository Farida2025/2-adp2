package usecase

import (
	"context"
	"fmt"
	"order-service/internal/domain"
	"order-service/internal/repository"

	"github.com/google/uuid"
)

type PaymentProvider interface {
	Authorize(ctx context.Context, orderID string, amount int64) (string, error)
}

type CreateOrder struct {
	repo            repository.OrderRepository
	paymentProvider PaymentProvider
}

func NewCreateOrder(repo repository.OrderRepository, pp PaymentProvider) *CreateOrder {
	return &CreateOrder{
		repo:            repo,
		paymentProvider: pp,
	}
}

type CreateOrderCommand struct {
	CustomerID string `json:"customer_id"`
	ItemName   string `json:"item_name"`
	Amount     int64  `json:"amount"`
}

func (uc *CreateOrder) Execute(ctx context.Context, cmd CreateOrderCommand) (string, error) {

	if cmd.Amount <= 0 {
		return "", fmt.Errorf("amount must be > 0")
	}

	order := domain.Order{
		ID:         uuid.New().String(),
		CustomerID: cmd.CustomerID,
		ItemName:   cmd.ItemName,
		Amount:     cmd.Amount,
		Status:     "Pending",
	}

	if err := uc.repo.Save(ctx, order); err != nil {
		return "", err
	}

	status, err := uc.paymentProvider.Authorize(ctx, order.ID, cmd.Amount)

	if err != nil {
		uc.repo.UpdateStatus(ctx, order.ID, "Failed")
		return "", fmt.Errorf("payment service unavailable: %w", err)
	}

	if status != "Authorized" {
		uc.repo.UpdateStatus(ctx, order.ID, "Failed")
		return "", fmt.Errorf("payment declined: status %s", status)
	}

	if err := uc.repo.UpdateStatus(ctx, order.ID, "Paid"); err != nil {
		return "", err
	}

	return order.ID, nil
}
