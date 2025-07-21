package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"transaction-management-system/consumer"
	"transaction-management-system/publisher"
	"transaction-management-system/transaction"
)

const (
	queueName = "casino"
	amqpURI   = "amqp://guest:guest@localhost:5672/"
)

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	// Wait for interrupt signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// Start publisher in a goroutine
	publisher := publisher.NewPublisher(amqpURI, queueName)
	wg.Add(1)
	go publisher.StartPublish(ctx, &wg, queueName)

	// Start consumer in goroutine
	consumer := consumer.NewConsumer(amqpURI, queueName)
	wg.Add(1)
	go consumer.Consume(ctx, &wg, queueName)

	// Start listening transaction api
	transactioApi := transaction.NewTransactionApi()
	wg.Add(1)
	go transactioApi.ListenAndServe(&wg)

	fmt.Println(" [*] Press CTRL+C to exit")
	<-sigChan
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()
}
