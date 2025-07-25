package publisher

import (
	"encoding/json"
	"sync"
	"testing"
	test "transaction-management-system/config"
	"transaction-management-system/transaction"

	"github.com/stretchr/testify/require"
)

func TestNewPublisher(t *testing.T) {
	t.Run("failed new publisher - wrong amqp uri", func(t *testing.T) {
		_, err := NewPublisher(test.WRONG_AMQP_URI, test.QUEUE_NAME)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to connect to RabbitMQ")
	})
	t.Run("failed new publisher - wrong queue name", func(t *testing.T) {
		_, err := NewPublisher(test.AMQP_URI, test.WRONG_QUEUE_NAME)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to declare a queue")
	})
	t.Run("succesfully created new publisher", func(t *testing.T) {
		p, err := NewPublisher(test.AMQP_URI, test.QUEUE_NAME)

		require.Nil(t, err)
		require.NotNil(t, p)
	})
}

func TestStartPublishing(t *testing.T) {
	var wg sync.WaitGroup
	t.Run("successfully published", func(t *testing.T) {
		p, _ := NewPublisher(test.AMQP_URI, test.QUEUE_NAME)

		wg.Add(1)
		go p.StartPublish(t.Context(), &wg, test.QUEUE_NAME)

		msgs, err := p.RabbitMQ.Consume(test.QUEUE_NAME)
		require.Nil(t, err)

		for msg := range msgs {
			var tr transaction.Transaction
			err = json.Unmarshal(msg.Body, &tr)

			require.Nil(t, err)
			require.NotEmpty(t, tr)
		}

		wg.Wait()
	})
}

func TestClose(t *testing.T) {
	t.Run("failed closing", func(t *testing.T) {
		p, _ := NewPublisher(test.AMQP_URI, test.QUEUE_NAME)
		p.Close()
		err := p.Close()

		require.Error(t, err)
		require.ErrorContains(t, err, "channel/connection is not open")
	})
	t.Run("succesfully closed", func(t *testing.T) {
		p, _ := NewPublisher(test.AMQP_URI, test.QUEUE_NAME)
		err := p.Close()

		require.Nil(t, err)

	})
}
