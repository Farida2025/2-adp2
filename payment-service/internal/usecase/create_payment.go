package usecase

import (
	"context"
	"payment-service/internal/domain"
	"payment-service/internal/repository"

	"github.com/google/uuid"
)

type CreatePayment struct {
	repo     repository.PaymentRepository
	notifier NotificationProvider
}

type NotificationProvider interface {
	SendPaymentNotification(ctx context.Context, orderID string, amount int64) error
}

func NewCreatePayment(repo repository.PaymentRepository, notifier NotificationProvider) *CreatePayment {
	return &CreatePayment{
		repo:     repo,
		notifier: notifier,
	}
}

type CreatePaymentCommand struct {
	OrderID string `json:"order_id"`
	Amount  int64  `json:"amount"`
}

type CreatePaymentResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

func (uc *CreatePayment) Execute(ctx context.Context, cmd CreatePaymentCommand) (CreatePaymentResponse, error) {
	if cmd.Amount > 100000 {
		return CreatePaymentResponse{Status: "Declined"}, nil
	}

	payment := domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       cmd.OrderID,
		TransactionID: uuid.New().String(),
		Amount:        cmd.Amount,
		Status:        "Authorized",
	}

	err := uc.repo.Save(ctx, payment)
	if err != nil {
		return CreatePaymentResponse{}, err
	}

	if uc.notifier != nil {
		err = uc.notifier.SendPaymentNotification(ctx, payment.OrderID, payment.Amount)
		if err != nil {
			return CreatePaymentResponse{}, err
		}
	}

	return CreatePaymentResponse{
		TransactionID: payment.TransactionID,
		Status:        "Authorized",
	}, nil
}
