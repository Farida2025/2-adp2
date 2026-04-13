package main

import (
	"database/sql"
	"log"
	"net"

	"payment-service/internal/app"
	grpcHandler "payment-service/internal/transport/grpc"

	pb "github.com/Farida2025/assignment2-generated/payment"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:0000@localhost:5432/payment_db?sslmode=disable")
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	application := app.NewApp(db)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen for gRPC: %v", err)
		}

		s := grpc.NewServer()

		paymentH := grpcHandler.NewPaymentGRPCHandler(application.PaymentUseCase)
		pb.RegisterPaymentServiceServer(s, paymentH)

		reflection.Register(s)

		log.Println("Payment gRPC Server is running on :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	r := gin.Default()
	application.Handler.RegisterRoutes(r)

	log.Println("Payment HTTP Service is running on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("failed to run HTTP: %v", err)
	}
}
