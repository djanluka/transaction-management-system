package main

import (
	"fmt"
	"os"
	"os/signal"
	"transaction-management-system/consumer"
	"transaction-management-system/publisher"
	"transaction-management-system/transaction"
)

const (
	queueName = "casino"
	amqpURI   = "amqp://guest:guest@localhost:5672/"
)

func main() {

	// Start publisher in a goroutine
	publisher := publisher.NewPublisher(amqpURI, queueName)
	// TODO: add wait group
	go publisher.StartPublish(queueName)

	// Start consumer in goroutine
	consumer := consumer.NewConsumer(amqpURI, queueName)
	// TODO: add wait group
	go consumer.Consume(queueName)

	// Start listening transaction api
	transactioApi := transaction.NewTransactionApi()
	// TODO: add wait group
	go transactioApi.ListenAndServe()

	// Wait for interrupt signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	fmt.Println(" [*] Press CTRL+C to exit")
	<-sigChan
	fmt.Println("Shutting down...")
}
