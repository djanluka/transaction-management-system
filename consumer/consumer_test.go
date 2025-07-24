package consumer

import (
	"testing"
	"transaction-management-system/test"

	"github.com/stretchr/testify/require"
)

func TestNewCustomer(t *testing.T) {

	t.Run("failed to connect to RabbitMQ", func(t *testing.T) {
		_, err := NewConsumer(test.WRONG_AMQP_URI, test.QUEUE_NAME)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to connect to RabbitMQ")
	})

	t.Run("failed to connect to DB", func(t *testing.T) {
		_, err := NewConsumer(test.AMQP_URI, test.QUEUE_NAME)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to load .env file")
	})

	t.Run("succesfully created consumer", func(t *testing.T) {
		c, err := NewConsumer(test.AMQP_URI, test.QUEUE_NAME)

		require.NoError(t, err)
		require.NotNil(t, c)
	})
}

func TestClose(t *testing.T) {

	t.Run("succesfully closed consumer", func(t *testing.T) {
		c, err := NewConsumer(test.AMQP_URI, test.WRONG_QUEUE_NAME)
		t.Log(err)
		t.Log(c)
		// c.Close()

		// _, err := c.Db.GetTransactions(t.Context(), nil, nil, 1)
		// require.Error(t, err)
		// require.ErrorContains(t, err, "failed to query transactions: sql: statement is closed")
	})

}
