package publisher

import (
	"fmt"
	"log"
	"time"
	"transaction-management-system/rabbitmq"
	"transaction-management-system/transaction"
)

type Publisher struct {
	RabbitMQ *rabbitmq.RabbitMQ
}

func NewPublisher(amqpURI string) *Publisher {
	rmq, err := rabbitmq.GetInstance(amqpURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	return &Publisher{
		RabbitMQ: rmq,
	}
}

func (p *Publisher) StartPublish(queueName string) {
	for i := 0; i < 5; i++ {
		transaction := transaction.NewTransaction()
		err := p.RabbitMQ.Publish(queueName, transaction)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
			continue
		}
		fmt.Printf(" [x] Sent: %s\n", transaction)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Publisher finished sending messages")
}
