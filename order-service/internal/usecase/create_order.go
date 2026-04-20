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

type OrderStatusNotifier interface {
	NotifyStatusChange(orderID, newStatus string)
}

type CreateOrder struct {
	repo            repository.OrderRepository
	paymentProvider PaymentProvider
	notifier        OrderStatusNotifier
}

func NewCreateOrder(repo repository.OrderRepository, pp PaymentProvider, n OrderStatusNotifier) *CreateOrder {
	return &CreateOrder{
		repo:            repo,
		paymentProvider: pp,
		notifier:        n,
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
	uc.notify(order.ID, "Pending")

	status, err := uc.paymentProvider.Authorize(ctx, order.ID, cmd.Amount)
	if err != nil {
		uc.repo.UpdateStatus(ctx, order.ID, "Failed")
		uc.notify(order.ID, "Failed")
		return "", fmt.Errorf("payment service unavailable: %w", err)
	}

	if status != "Authorized" {
		uc.repo.UpdateStatus(ctx, order.ID, "Failed")
		uc.notify(order.ID, "Failed")
		return "", fmt.Errorf("payment declined: status %s", status)
	}

	if err := uc.repo.UpdateStatus(ctx, order.ID, "Paid"); err != nil {
		return "", err
	}

	uc.notify(order.ID, "Paid")
	return order.ID, nil
}

func (uc *CreateOrder) notify(id, status string) {
	if uc.notifier != nil {
		uc.notifier.NotifyStatusChange(id, status)
	}
}
