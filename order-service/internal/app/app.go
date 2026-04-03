package app

import (
	"database/sql"
	"order-service/internal/repository/postgres"
	"order-service/internal/transport/http"
	"order-service/internal/usecase"
)

type App struct {
	Handler *http.Handler
}

func NewApp(db *sql.DB, paymentURL string) *App {
	repo := postgres.NewOrderRepo(db)
	Uc := usecase.NewCreateOrder(repo, paymentURL)
	recentUc := usecase.NewGetRecentOrders(repo)
	handler := http.NewHandler(Uc, recentUc)
	return &App{Handler: handler}
}
