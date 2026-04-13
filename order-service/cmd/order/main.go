package main

import (
	"database/sql"
	"log"
	"order-service/internal/app"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:0000@localhost:5432/order_db?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	paymentGRPCAddr := os.Getenv("PAYMENT_GRPC_ADDR")
	if paymentGRPCAddr == "" {
		log.Fatal("PAYMENT_GRPC_ADDR is not set")
	}

	application := app.NewApp(db, paymentGRPCAddr)
	defer application.GRPCAdapter.Close()

	r := gin.Default()
	application.Handler.RegisterRoutes(r)

	log.Println("Order Service running on :8080 (gRPC enabled)")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
