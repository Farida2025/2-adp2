package grpc

import (
	"context"
	"payment-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Farida2025/assignment2-generated/payment"
)

type PaymentGRPCHandler struct {
	pb.UnimplementedPaymentServiceServer
	uc     *usecase.CreatePayment
}

func NewPaymentGRPCHandler(uc *usecase.CreatePayment) *PaymentGRPCHandler {
	return &PaymentGRPCHandler{uc: uc}

}

func (h *PaymentGRPCHandler) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {

	resp, err := h.uc.Execute(ctx, usecase.CreatePaymentCommand{
		OrderID: req.GetOrderId(),
		Amount:  req.GetAmount(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if resp.Status == "Declined" {
		return nil, status.Errorf(codes.FailedPrecondition, "payment declined")
	}

	return &pb.PaymentResponse{
		TransactionId: resp.TransactionID,
		Status:        resp.Status,
	}, nil
}

