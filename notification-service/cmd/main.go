package main

import (
	"notification-service/internal/consumer"
	"os"
)

func main() {
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}
	consumer.Start(amqpURL)
}
