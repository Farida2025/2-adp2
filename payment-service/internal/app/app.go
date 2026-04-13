package app

import (
	"database/sql"
	"payment-service/internal/repository/postgres"
	"payment-service/internal/transport/http"
	"payment-service/internal/usecase"
)

type App struct {
	Handler        *http.Handler
	PaymentUseCase *usecase.CreatePayment
}

func NewApp(db *sql.DB) *App {
	repo := postgres.NewPaymentRepo(db)
	uc := usecase.NewCreatePayment(repo)
	handler := http.NewHandler(uc)

	return &App{
		Handler:        handler,
		PaymentUseCase: uc,
	}
}
