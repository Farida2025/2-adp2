package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPAdapter struct {
	client     *http.Client
	paymentURL string
}

func NewHTTPAdapter(client *http.Client, url string) *HTTPAdapter {
	return &HTTPAdapter{
		client:     client,
		paymentURL: url,
	}
}

func (a *HTTPAdapter) Authorize(ctx context.Context, orderID string, amount int64) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"order_id": orderID,
		"amount":   amount,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal payment request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.paymentURL+"/payments", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("payment service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "Declined", nil
	}

	var result struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Status, nil
}
