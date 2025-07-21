package rabbitmq

import (
	"math"
	"testing"
	"transaction-management-system/transaction"

	"github.com/stretchr/testify/require"
)

func TestGetInstance(t *testing.T) {
	amqpURI := "amqp://guest:guest@localhost:5672/"
	wAmqpURI := "amqp://wrongUri"
	queueName := "casino"
	wQueueName := "amq.queueName"

	t.Run("wrong amqpUri", func(t *testing.T) {
		_, err := GetInstance(wAmqpURI, wQueueName)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to connect to RabbitMQ")
	})

	t.Run("wrong queue name", func(t *testing.T) {
		_, err := GetInstance(amqpURI, wQueueName)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to declare a queue")
	})

	t.Run("succesfull connection to RabbitMQ", func(t *testing.T) {
		rmq, err := GetInstance(amqpURI, queueName)
		require.NoError(t, err)
		require.NotNil(t, rmq)
	})
}

func TestClose(t *testing.T) {
	amqpURI := "amqp://guest:guest@localhost:5672/"
	queueName := "casino"

	rmq, _ := GetInstance(amqpURI, queueName)
	t.Run("succesfully closed rmq", func(t *testing.T) {
		rmq.Close()

		err := rmq.Publish(queueName, transaction.NewTransaction())
		require.Error(t, err)
		require.ErrorContains(t, err, "channel/connection is not open")
	})
}

func TestPublish(t *testing.T) {
	amqpURI := "amqp://guest:guest@localhost:5672/"
	queueName := "casino"

	rmq, _ := GetInstance(amqpURI, queueName)
	t.Run("failed unmarshaling transaction", func(t *testing.T) {
		tr := transaction.NewTransaction()
		tr.Amount = math.NaN()

		err := rmq.Publish(queueName, tr)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to marshal transaction")
	})

	t.Run("succesfully publishing", func(t *testing.T) {
		err := rmq.Publish(queueName, transaction.NewTransaction())
		require.NoError(t, err)
		require.Nil(t, err)
	})

	t.Run("failed publishing when rmq is closed", func(t *testing.T) {
		rmq.Close()

		err := rmq.Publish(queueName, transaction.NewTransaction())
		require.Error(t, err)
		require.ErrorContains(t, err, "channel/connection is not open")
	})
}

func TestConsume(t *testing.T) {
	amqpURI := "amqp://guest:guest@localhost:5672/"
	queueName := "casino"
	wQueueName := "amq.queueName"

	t.Run("failed set Qos", func(t *testing.T) {
		rmq, _ := GetInstance(amqpURI, queueName)
		rmq.Close()

		_, err := rmq.Consume(queueName)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to set Qos")
	})

	t.Run("failed to read from wrong queue", func(t *testing.T) {
		rmq, _ := GetInstance(amqpURI, queueName)
		defer rmq.Close()

		_, err := rmq.Consume(wQueueName)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to register a consumer")
	})
	t.Run("succesfully consume", func(t *testing.T) {
		rmq, _ := GetInstance(amqpURI, queueName)
		defer rmq.Close()

		_, err := rmq.Consume(queueName)
		require.NoError(t, err)
	})

}
