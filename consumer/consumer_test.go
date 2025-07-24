package consumer

import (
	"bytes"
	"context"
	"log"
	"os"
	"sync"
	"testing"
	"time"
	"transaction-management-system/test"
	"transaction-management-system/transaction"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Load .env from project root
	_ = godotenv.Load("../.env")

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestNewConsumer(t *testing.T) {
	t.Run("failed new consumer - wrong amqp uri", func(t *testing.T) {
		_, err := NewConsumer(test.WRONG_AMQP_URI, test.QUEUE_NAME)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to connect to RabbitMQ")
	})
	t.Run("failed new consumer - wrong queue name", func(t *testing.T) {
		_, err := NewConsumer(test.AMQP_URI, test.WRONG_QUEUE_NAME)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to declare a queue")
	})
	t.Run("succesfully created new consumer", func(t *testing.T) {
		c, err := NewConsumer(test.AMQP_URI, test.QUEUE_NAME)

		require.Nil(t, err)
		require.NotNil(t, c)
	})
}

func TestConsume(t *testing.T) {

	var wg sync.WaitGroup
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		// Reset log output when test is done
		log.SetOutput(os.Stderr)
	}()

	t.Run("failed consuming", func(t *testing.T) {
		c, _ := NewConsumer(test.AMQP_URI, test.QUEUE_NAME)
		defer c.Close()

		c.RabbitMQ.Close()

		wg.Add(1)
		go c.Consume(t.Context(), &wg, test.QUEUE_NAME)
		wg.Wait()

		logOutput := buf.String()
		require.Contains(t, logOutput, "Failed to start consumer:",
			"Expected error log not found in:\n%s", logOutput)
	})

	t.Run("successfully consumed", func(t *testing.T) {
		c, _ := NewConsumer(test.AMQP_URI, test.QUEUE_NAME)
		defer c.Close()

		c.RabbitMQ.Publish(test.QUEUE_NAME, transaction.NewTransaction())

		ctx, close := context.WithTimeout(t.Context(), 1*time.Second)
		defer close()

		wg.Add(1)
		go c.Consume(ctx, &wg, test.QUEUE_NAME)

		wg.Wait()
	})

	t.Run("successfully consumed", func(t *testing.T) {
		c, _ := NewConsumer(test.AMQP_URI, test.QUEUE_NAME)
		defer c.Close()

		c.RabbitMQ.Publish(test.QUEUE_NAME, transaction.NewTransaction())

		ctx, close := context.WithTimeout(t.Context(), 1*time.Second)
		defer close()

		wg.Add(1)
		go c.Consume(ctx, &wg, test.QUEUE_NAME)

		wg.Wait()

		logOutput := buf.String()
		require.Contains(t, logOutput, "Received",
			"Expected error log not found in:\n%s", logOutput)
	})

}

func TestClose(t *testing.T) {
	t.Run("failed closing", func(t *testing.T) {
		c, _ := NewConsumer(test.AMQP_URI, test.QUEUE_NAME)
		c.Close()
		err := c.Close()

		require.Error(t, err)
		require.ErrorContains(t, err, "channel/connection is not open")
	})
	t.Run("succesfully closed", func(t *testing.T) {
		c, _ := NewConsumer(test.AMQP_URI, test.QUEUE_NAME)
		err := c.Close()

		require.Nil(t, err)
	})
}
