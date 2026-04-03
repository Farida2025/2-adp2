package main

import (
	"database/sql"
	"log"
	"order-service/internal/app"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:0000@localhost:5432/order_db?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	paymentURL := "http://localhost:8081"

	app := app.NewApp(db, paymentURL)

	r := gin.Default()
	app.Handler.RegisterRoutes(r)

	log.Println("Order Service running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
