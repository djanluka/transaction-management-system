package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"transaction-management-system/config"
	"transaction-management-system/database"
	"transaction-management-system/rabbitmq"
	"transaction-management-system/transaction"
)

type Consumer struct {
	RabbitMQ *rabbitmq.RabbitMQ
	Db       *database.Database
}

func NewConsumer(amqpURI, queueName string) (*Consumer, error) {
	rmq, err := rabbitmq.GetInstance(amqpURI, queueName)
	if err != nil {
		return nil, err
	}

	db, err := database.GetDB(config.DB_SCHEMA)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		RabbitMQ: rmq,
		Db:       db,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context, wg *sync.WaitGroup, queueName string) {
	defer wg.Done()
	defer c.Close()

	msgs, err := c.RabbitMQ.Consume(queueName)
	if err != nil {
		log.Printf("Failed to start consumer: %v\n", err)
		return
	}

	fmt.Println("Consumer started. Waiting for messages...")
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			// Unmarshal body as transaction
			var tr transaction.Transaction
			err := json.Unmarshal(msg.Body, &tr)
			if err != nil {
				log.Printf("Error decoding transaction: %s\n", err)
				continue
			}

			// Insert transaction into database
			log.Printf(" [x] Received: %s\n", tr)
			if err := c.Db.InsertTransaction(tr.UserId, tr.TransactionType, tr.Amount, tr.Timestamp); err != nil {
				log.Printf(" WARN: Message has not been processed successfully")
				msg.Nack(false, true)
				continue
			}
			log.Printf(" [x] Inserted: %s\n", tr)
		}
	}
}

func (c *Consumer) Close() error {
	if err := c.RabbitMQ.Close(); err != nil {
		return err
	}
	if err := c.Db.Close(); err != nil {
		return err
	}
	log.Println("Consumer closed succesfully")
	return nil
}
