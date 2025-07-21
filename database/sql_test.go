package database

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBInitializationFailure(t *testing.T) {
	// 1. Backup current state
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// 2. Create and move to a temp directory without .env
	tempDir := t.TempDir()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	t.Run("succesfull no finding.env file", func(t *testing.T) {
		_, err := GetDB()
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to load .env file")

	})
}

// TestGetDBSingleton tests the singleton behavior
func TestGetDBSingleton(t *testing.T) {
	// Reset the singleton for testing
	instance = nil
	once = sync.Once{}

	t.Run("singleton pattern", func(t *testing.T) {
		// First call should succeed
		db1, err := GetDB()
		require.NoError(t, err)
		require.NotNil(t, db1)

		// Second call should return the same instance
		db2, err := GetDB()
		require.NoError(t, err)
		require.NotNil(t, db2)

		// Verify it's the same instance
		assert.Equal(t, db1, db2)
	})
}

// TestInsertTransaction tests the InsertTransaction method
func TestInsertTransaction(t *testing.T) {

	db, _ := GetDB()
	userId := 999
	transactionType := "bet"
	wrongTransactionType := "wrong_type"
	amount := 0.5

	// Test cases
	t.Run("successful insert with valid data", func(t *testing.T) {
		err := db.InsertTransaction(userId, transactionType, amount, time.Now())
		require.NoError(t, err)
	})

	t.Run("error with invalid transaction type", func(t *testing.T) {
		err := db.InsertTransaction(userId, wrongTransactionType, amount, time.Now())
		require.Error(t, err)
		require.ErrorContains(t, err, "Data truncated for column 'transaction_type'")
	})
}

// TestGetTransactions tests the GetTransactions method
func TestGetTransactions(t *testing.T) {

	db, _ := GetDB()
	userId := "1"
	transactionType := "bet"

	// Test cases
	t.Run("successful get transaction", func(t *testing.T) {
		rows, err := db.GetTransactions(t.Context(), &userId, &transactionType, 1)
		require.NoError(t, err)
		defer rows.Close()

		assert.Equal(t, rows.Next(), true)
	})
}

func TestClose(t *testing.T) {
	db, _ := GetDB()
	userId := "1"
	transactionType := "bet"

	t.Run("successful close", func(t *testing.T) {
		db.Close()

		_, err := db.GetTransactions(t.Context(), &userId, &transactionType, 1)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to query transactions: sql: statement is closed")
	})

}
