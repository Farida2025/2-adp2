package app

import (
	"database/sql"
	"log"
	"payment-service/internal/infrastructure/rabbitmq"
	"payment-service/internal/repository/postgres"
	"payment-service/internal/transport/http"
	"payment-service/internal/usecase"
)

type App struct {
	Handler        *http.Handler
	PaymentUseCase *usecase.CreatePayment
	RabbitMQ       *rabbitmq.RabbitMQProducer
}

func NewApp(db *sql.DB, amqpURL string) *App {
	repo := postgres.NewPaymentRepo(db)

	notifier, err := rabbitmq.NewRabbitMQProducer(amqpURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	uc := usecase.NewCreatePayment(repo, notifier)
	handler := http.NewHandler(uc)

	return &App{
		Handler:        handler,
		PaymentUseCase: uc,
		RabbitMQ:       notifier,
	}
}
