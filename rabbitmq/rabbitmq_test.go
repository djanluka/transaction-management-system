package rabbitmq

import (
	"math"
	"testing"
	"transaction-management-system/test"
	"transaction-management-system/transaction"

	"github.com/stretchr/testify/require"
)

func TestGetInstance(t *testing.T) {

	t.Run("wrong amqp uri", func(t *testing.T) {
		_, err := GetInstance(test.WRONG_AMQP_URI, test.WRONG_QUEUE_NAME)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to connect to RabbitMQ")
	})

	t.Run("wrong queue name", func(t *testing.T) {
		_, err := GetInstance(test.AMQP_URI, test.WRONG_QUEUE_NAME)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to declare a queue")
	})

	t.Run("succesfull connection", func(t *testing.T) {
		rmq, err := GetInstance(test.AMQP_URI, test.QUEUE_NAME)
		require.NoError(t, err)
		require.NotNil(t, rmq)
	})
}

func TestClose(t *testing.T) {
	rmq, _ := GetInstance(test.AMQP_URI, test.QUEUE_NAME)

	t.Run("failed closing rmq when connection is already closed", func(t *testing.T) {
		rmq.conn.Close()
		err := rmq.Close()

		require.Error(t, err)
		require.ErrorContains(t, err, "channel/connection is not open")
	})

	t.Run("succesfully closed rmq", func(t *testing.T) {
		rmq.Close()

		err := rmq.Publish(test.QUEUE_NAME, transaction.NewTransaction())
		require.Error(t, err)
		require.ErrorContains(t, err, "channel/connection is not open")
	})
}

func TestPublish(t *testing.T) {
	t.Run("failed unmarshaling transaction", func(t *testing.T) {
		rmq, _ := GetInstance(test.AMQP_URI, test.QUEUE_NAME)
		defer rmq.Close()
		tr := transaction.NewTransaction()
		tr.Amount = math.NaN()

		err := rmq.Publish(test.QUEUE_NAME, tr)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to marshal transaction")
	})

	t.Run("succesfully publishing", func(t *testing.T) {
		rmq, _ := GetInstance(test.AMQP_URI, test.QUEUE_NAME)
		defer rmq.Close()

		err := rmq.Publish(test.QUEUE_NAME, transaction.NewTransaction())
		require.NoError(t, err)
		require.Nil(t, err)
	})

	t.Run("failed publishing when rmq is closed", func(t *testing.T) {
		rmq, _ := GetInstance(test.AMQP_URI, test.QUEUE_NAME)
		rmq.Close()

		err := rmq.Publish(test.QUEUE_NAME, transaction.NewTransaction())
		require.Error(t, err)
		require.ErrorContains(t, err, "channel/connection is not open")
	})
}

func TestConsume(t *testing.T) {
	t.Run("failed set Qos", func(t *testing.T) {
		rmq, _ := GetInstance(test.AMQP_URI, test.QUEUE_NAME)
		rmq.Close()

		_, err := rmq.Consume(test.QUEUE_NAME)
		require.Error(t, err)
		require.ErrorContains(t, err, "channel/connection is not open")
	})

	t.Run("failed to read from wrong queue", func(t *testing.T) {
		rmq, _ := GetInstance(test.AMQP_URI, test.QUEUE_NAME)
		defer rmq.Close()

		_, err := rmq.Consume(test.WRONG_QUEUE_NAME)
		require.Error(t, err)
		require.ErrorContains(t, err, "NOT_FOUND")
	})
	t.Run("succesfully consume", func(t *testing.T) {
		rmq, _ := GetInstance(test.AMQP_URI, test.QUEUE_NAME)
		defer rmq.Close()

		_, err := rmq.Consume(test.QUEUE_NAME)
		require.NoError(t, err)
	})

}
