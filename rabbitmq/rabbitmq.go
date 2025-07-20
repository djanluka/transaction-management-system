package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"transaction-management-system/transaction"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ represents a wrapper for RabbitMQ connection and channel
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// GetInstance returns a instance of RabbitMQ
func GetInstance(amqpURI, queueName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, errors.New("failed to connect to RabbitMQ")
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, errors.New("failed to open a channel")
	}

	_, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, errors.New("failed to declare a queue")
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQ) Close() error {

	var err error
	if r.channel != nil {
		if closeErr := r.channel.Close(); closeErr != nil {
			err = fmt.Errorf("channel close error: %v", closeErr)
		}
	}

	if r.conn != nil {
		if closeErr := r.conn.Close(); closeErr != nil {
			if err != nil {
				err = fmt.Errorf("%v, connection close error: %v", err, closeErr)
			} else {
				err = fmt.Errorf("connection close error: %v", closeErr)
			}
		}
	}

	return err
}

func (r *RabbitMQ) Publish(queueName string, transaction transaction.Transaction) error {
	// Marshal to JSON
	body, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	err = r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	return nil
}

func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	if err := r.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return nil, errors.New("failed to set Qos")
	}

	msgs, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %v", err)
	}

	return msgs, nil
}
