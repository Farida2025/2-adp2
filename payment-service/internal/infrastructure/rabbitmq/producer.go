package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"
)

type RabbitMQProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQProducer(url string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"payment.completed",
		true,
		false,
		false,
		false,
		nil,
	)

	return &RabbitMQProducer{conn: conn, channel: ch}, nil
}

func (p *RabbitMQProducer) SendPaymentNotification(ctx context.Context, orderID string, amount int64) error {
	payload := map[string]interface{}{
		"order_id":       orderID,
		"amount":         float64(amount) / 100.0,
		"customer_email": "user@example.com",
		"status":         "success",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",
		"payment.completed",
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (p *RabbitMQProducer) Close() {
	p.channel.Close()
	p.conn.Close()
}
