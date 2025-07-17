package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"transaction-management-system/rabbitmq"
)

const (
	queueName = "hello"
	amqpURI   = "amqp://guest:guest@localhost:5672/"
)

func main() {
	// Initialize RabbitMQ singleton
	rmq, err := rabbitmq.GetInstance("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// Start consumer in a goroutine
	go func() {
		msgs, err := rmq.Consume(queueName)
		if err != nil {
			log.Printf("Failed to start consumer: %v", err)
			return
		}

		fmt.Println("Consumer started. Waiting for messages...")
		for msg := range msgs {
			fmt.Printf(" [x] Received: %s\n", msg.Body)
		}
	}()

	// Start publisher in a goroutine
	go func() {
		for i := 0; i < 5; i++ {
			message := fmt.Sprintf("Hello World %d", i)
			err := rmq.Publish(queueName, message)
			if err != nil {
				log.Printf("Failed to publish message: %v", err)
				continue
			}
			fmt.Printf(" [x] Sent: %s\n", message)
			time.Sleep(1 * time.Second)
		}
		fmt.Println("Publisher finished sending messages")
	}()

	// Wait for interrupt signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	fmt.Println(" [*] Press CTRL+C to exit")
	<-sigChan
	fmt.Println("Shutting down...")
}
