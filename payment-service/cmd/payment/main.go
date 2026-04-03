package main

import (
	"database/sql"
	"log"
	"payment-service/internal/app"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:0000@localhost:5432/payment_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := app.NewApp(db)

	r := gin.Default()
	app.Handler.RegisterRoutes(r)

	log.Println("Payment Service running on :8081")
	r.Run(":8081")
}
