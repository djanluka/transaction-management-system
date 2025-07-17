package rabbitmq

import (
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ represents a wrapper for RabbitMQ connection and channel
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// Add other fields you might need (queues, exchanges, etc.)
}

var (
	instance *RabbitMQ
	once     sync.Once
)

// GetInstance returns a singleton instance of RabbitMQ
func GetInstance(amqpURI string) (*RabbitMQ, error) {
	var initErr error

	once.Do(func() {
		conn, err := amqp.Dial(amqpURI)
		if err != nil {
			initErr = fmt.Errorf("failed to connect to RabbitMQ: %v", err)
			return
		}

		channel, err := conn.Channel()
		if err != nil {
			initErr = fmt.Errorf("failed to open a channel: %v", err)
			return
		}

		instance = &RabbitMQ{
			conn:    conn,
			channel: channel,
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return instance, nil
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

	instance = nil // Reset instance to allow reconnection if needed
	return err
}

// GetChannel returns the AMQP channel
func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.channel
}

// GetConnection returns the AMQP connection
func (r *RabbitMQ) GetConnection() *amqp.Connection {
	return r.conn
}

func (r *RabbitMQ) Publish(queueName string, body string) error {
	_, err := r.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	err = r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	return nil
}

func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	_, err := r.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
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
