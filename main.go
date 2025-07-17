package main

import (
	"fmt"
	"os"
	"os/signal"
	"transaction-management-system/consumer"
	"transaction-management-system/publisher"
)

const (
	queueName = "hello"
	amqpURI   = "amqp://guest:guest@localhost:5672/"
)

func main() {
	// Start publisher in a goroutine
	publisher := publisher.NewPublisher(amqpURI)
	// TODO: add wait group
	go publisher.StartPublish(queueName)

	// Start consumer in goroutine
	consumer := consumer.NewConsumer(amqpURI)
	// TODO: add wait group
	go consumer.Consume(queueName)

	// Wait for interrupt signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	fmt.Println(" [*] Press CTRL+C to exit")
	<-sigChan
	fmt.Println("Shutting down...")
}
