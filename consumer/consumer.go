package consumer

import (
	"fmt"
	"log"
	"transaction-management-system/rabbitmq"
)

type Consumer struct {
	RabbitMQ *rabbitmq.RabbitMQ
}

func NewConsumer(amqpURI string) *Consumer {
	rmq, err := rabbitmq.GetInstance(amqpURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	return &Consumer{
		RabbitMQ: rmq,
	}
}

func (c *Consumer) Consume(queueName string) {
	msgs, err := c.RabbitMQ.Consume(queueName)
	if err != nil {
		log.Printf("Failed to start consumer: %v", err)
		return
	}

	fmt.Println("Consumer started. Waiting for messages...")
	for msg := range msgs {
		fmt.Printf(" [x] Received: %s\n", msg.Body)
	}
}
