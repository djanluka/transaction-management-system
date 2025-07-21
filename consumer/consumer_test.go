package consumer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCustomer(t *testing.T) {
	amqpURI := "amqp://guest:guest@localhost:5672/"
	wAmqpURI := "amqp://wrongUri"
	queueName := "casino"

	t.Run("failed to connect to RabbitMQ", func(t *testing.T) {
		_, err := NewConsumer(wAmqpURI, queueName)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to connect to RabbitMQ")
	})

	t.Run("failed to connect to DB", func(t *testing.T) {
		_, err := NewConsumer(amqpURI, queueName)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to load .env file")
	})

	t.Run("succesfully created consumer", func(t *testing.T) {
		c, err := NewConsumer(amqpURI, queueName)

		require.NoError(t, err)
		require.NotNil(t, c)
	})
}

func TestClose(t *testing.T) {
	amqpURI := "amqp://guest:guest@localhost:5672/"
	// wAmqpURI := "amqp://wrongUri"
	queueName := "casino"

	t.Run("succesfully closed consumer", func(t *testing.T) {
		c, err := NewConsumer(amqpURI, queueName)
		t.Log(err)
		t.Log(c)
		// c.Close()

		// _, err := c.Db.GetTransactions(t.Context(), nil, nil, 1)
		// require.Error(t, err)
		// require.ErrorContains(t, err, "failed to query transactions: sql: statement is closed")
	})

}
