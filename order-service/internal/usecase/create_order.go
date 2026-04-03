package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/internal/domain"
	"order-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type CreateOrder struct {
	repo          repository.OrderRepository
	paymentClient *http.Client
	paymentURL    string
}

func NewCreateOrder(repo repository.OrderRepository, paymentURL string) *CreateOrder {
	return &CreateOrder{
		repo:          repo,
		paymentClient: &http.Client{Timeout: 2 * time.Second},
		paymentURL:    paymentURL,
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

	payload := map[string]interface{}{
		"order_id": order.ID,
		"amount":   cmd.Amount,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := uc.paymentClient.Post(
		uc.paymentURL+"/payments",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		uc.repo.UpdateStatus(ctx, order.ID, "Failed")
		return "", fmt.Errorf("payment service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		uc.repo.UpdateStatus(ctx, order.ID, "Failed")
		return "", fmt.Errorf("payment declined")
	}

	uc.repo.UpdateStatus(ctx, order.ID, "Paid")
	return order.ID, nil
}
