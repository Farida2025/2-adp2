package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/streadway/amqp"
)

type PaymentEvent struct {
	OrderID string  `json:"order_id"`
	Email   string  `json:"customer_email"`
	Amount  float64 `json:"amount"`
}

var (
	processedOrders = make(map[string]bool)
	mu              sync.Mutex
)

func Start(amqpURL string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"payment.dlx",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare DLX: %v", err)
	}

	dlq, err := ch.QueueDeclare(
		"payment.failed",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare DLQ: %v", err)
	}

	err = ch.QueueBind(dlq.Name, "payment.failed", "payment.dlx", false, nil)
	if err != nil {
		log.Fatalf("Failed to bind DLQ: %v", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange":    "payment.dlx",
		"x-dead-letter-routing-key": "payment.failed",
	}

	q, err := ch.QueueDeclare(
		"payment.completed",
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		log.Fatalf("Queue declare failed: %v. Hint: if queue already exists without DLQ, delete it first.", err)
	}

	err = ch.Qos(1, 0, false)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Consume failed: %v", err)
	}

	log.Println("[*] Notification Service started. Waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down gracefully...")
			return

		case d, ok := <-msgs:
			if !ok {
				log.Println("Channel closed, exiting...")
				return
			}

			var event PaymentEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error decoding message: %v. Sending to DLQ.", err)
				d.Nack(false, false)
				continue
			}

			if event.Amount == 666.66 {
				log.Printf("[DLQ] Permanent error for Order #%s. Moving to payment.failed", event.OrderID)
				d.Nack(false, false)
				continue
			}

			mu.Lock()
			if _, exists := processedOrders[event.OrderID]; exists {
				log.Printf("[Idempotency] Order #%s already processed. Skipping...", event.OrderID)
				mu.Unlock()
				d.Ack(false)
				continue
			}

			log.Printf("[Notification] Sent email to %s for Order #%s. Amount: $%.2f",
				event.Email, event.OrderID, event.Amount)

			processedOrders[event.OrderID] = true
			mu.Unlock()
			d.Ack(false)
		}
	}
}
