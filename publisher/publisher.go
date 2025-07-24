package publisher

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"transaction-management-system/rabbitmq"
	"transaction-management-system/transaction"
)

type Publisher struct {
	RabbitMQ *rabbitmq.RabbitMQ
}

func NewPublisher(amqpURI, queueName string) (*Publisher, error) {
	rmq, err := rabbitmq.GetInstance(amqpURI, queueName)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		RabbitMQ: rmq,
	}, nil
}

func (p *Publisher) StartPublish(ctx context.Context, wg *sync.WaitGroup, queueName string) {
	defer wg.Done()
	defer p.Close()

	publishingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case <-publishingCtx.Done():
			return
		default:
			transaction := transaction.NewTransaction()
			err := p.RabbitMQ.Publish(queueName, transaction)
			if err != nil {
				log.Printf("Failed to publish message: %v", err)
				continue
			}
			fmt.Printf(" [x] Sent: %s\n", transaction)
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *Publisher) Close() error {
	if err := p.RabbitMQ.Close(); err != nil {
		return err
	}
	log.Println("Publisher closed succesfully")
	return nil
}
