package publisher

import (
	"fmt"
	"log"
	"time"
	"transaction-management-system/rabbitmq"
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
		message := fmt.Sprintf("Hello World %d", i)
		err := p.RabbitMQ.Publish(queueName, message)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
			continue
		}
		fmt.Printf(" [x] Sent: %s\n", message)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Publisher finished sending messages")
}
