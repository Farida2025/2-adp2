package payment

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Farida2025/assignment2-generated/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCAdapter struct {
	client pb.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCAdapter(address string) (*GRPCAdapter, error) {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}

	return &GRPCAdapter{
		client: pb.NewPaymentServiceClient(conn),
		conn:   conn,
	}, nil
}

func (a *GRPCAdapter) Close() error {
	return a.conn.Close()
}

func (a *GRPCAdapter) Authorize(ctx context.Context, orderID string, amount int64) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := a.client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	})
	if err != nil {
		return "", fmt.Errorf("grpc payment failed: %w", err)
	}

	return resp.GetStatus(), nil
}
