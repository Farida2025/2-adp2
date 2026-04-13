package app

import (
	"database/sql"
	"log"
	"order-service/internal/infrastructure/payment"
	"order-service/internal/repository/postgres"
	delivery "order-service/internal/transport/http"
	"order-service/internal/usecase"
)

type App struct {
	Handler     *delivery.Handler
	GRPCAdapter *payment.GRPCAdapter
}

func NewApp(db *sql.DB, paymentGRPCAddr string) *App {
	repo := postgres.NewOrderRepo(db)

	paymentAdapter, err := payment.NewGRPCAdapter(paymentGRPCAddr)
	if err != nil {
		log.Fatal("Failed to create payment gRPC adapter:", err)
	}

	getUc := usecase.NewGetOrder(repo)
	cancelUc := usecase.NewCancelOrder(repo)

	createOrderUc := usecase.NewCreateOrder(repo, paymentAdapter)

	h := delivery.NewHandler(createOrderUc, getUc, cancelUc)

	return &App{
		Handler:     h,
		GRPCAdapter: paymentAdapter,
	}
}
