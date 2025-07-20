package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"transaction-management-system/database"
	"transaction-management-system/rabbitmq"
	"transaction-management-system/transaction"
)

type Consumer struct {
	RabbitMQ *rabbitmq.RabbitMQ
	Db       *database.Database
}

func NewConsumer(amqpURI, queueName string) *Consumer {
	rmq, err := rabbitmq.GetInstance(amqpURI, queueName)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	db, err := database.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to Database: %v", err)
	}

	return &Consumer{
		RabbitMQ: rmq,
		Db:       db,
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
		// Unmarshal body as tr
		var tr transaction.Transaction
		err := json.Unmarshal(msg.Body, &tr)
		if err != nil {
			log.Printf("Error decoding transaction: %s", err)
			continue
		}

		// Insert transaction into database
		fmt.Printf(" [x] Received: %s\n", tr)
		if err := c.Db.InsertTransaction(tr.UserId, tr.TransactionType, tr.Amount, tr.Timestamp); err != nil {
			msg.Nack(false, true)
			continue
		}

		fmt.Printf(" [x] Inserted: %s\n", tr)
	}
}
