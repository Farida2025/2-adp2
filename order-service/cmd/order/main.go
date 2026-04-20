package main

import (
	"database/sql"
	"log"
	"net"
	"order-service/internal/app"
	"os"

	"google.golang.org/grpc/reflection"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	grpcHandler "order-service/internal/transport/grpc"

	pb "github.com/Farida2025/assignment2-generated/order"
)

func main() {
	godotenv.Load()

	db, err := sql.Open("postgres", "postgres://postgres:0000@localhost:5432/order_db?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	paymentGRPCAddr := os.Getenv("PAYMENT_GRPC_ADDR")
	if paymentGRPCAddr == "" {
		log.Fatal("PAYMENT_GRPC_ADDR is not set")
	}

	streamSrv := grpcHandler.NewOrderStreamServer()
	application := app.NewApp(db, paymentGRPCAddr, streamSrv)
	defer application.GRPCAdapter.Close()

	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("Failed to listen for gRPC: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterOrderServiceServer(s, streamSrv)
		reflection.Register(s)

		log.Println("gRPC Streaming Server running on :50052")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	r := gin.Default()
	application.Handler.RegisterRoutes(r)

	log.Println("Order Service (HTTP) running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
