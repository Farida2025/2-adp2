package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"payment-service/internal/app"
	grpcHandler "payment-service/internal/transport/grpc"

	pb "github.com/Farida2025/assignment2-generated/payment"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("Method: %s, Duration: %s, Error: %v", info.FullMethod, time.Since(start), err)
	return resp, err
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	application := app.NewApp(db, amqpURL)
	defer application.RabbitMQ.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen for gRPC: %v", err)
		}
		s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
		paymentH := grpcHandler.NewPaymentGRPCHandler(application.PaymentUseCase)
		pb.RegisterPaymentServiceServer(s, paymentH)
		reflection.Register(s)

		log.Println("Payment gRPC Server is running on :50051")
		if err := s.Serve(lis); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	go func() {
		r := gin.Default()
		application.Handler.RegisterRoutes(r)
		log.Println("Payment HTTP Service is running on :8081")
		if err := r.Run(":8081"); err != nil {
			log.Printf("HTTP server stopped: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down payment-service gracefully...")
}
