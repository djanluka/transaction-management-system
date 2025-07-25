package database

import (
	"os"
	"testing"
	"time"
	test "transaction-management-system/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBInitializationFailure(t *testing.T) {
	// Backup current state
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create and move to a temp directory without .env
	tempDir := t.TempDir()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	t.Run("succesfull no finding.env file", func(t *testing.T) {
		_, err := GetDB(test.DB_SCHEMA)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to load .env file")

	})
}

// TestGetDBSingleton tests the singleton behavior
func TestGetDBSingleton(t *testing.T) {
	t.Run("failed ENV_PATH ", func(t *testing.T) {
		ResetInstance()

		os.Setenv("ENV_PATH", "")
		_, err := GetDB(test.DB_SCHEMA)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to load .env file")

		os.Setenv("ENV_PATH", "../.env")
	})

	t.Run("failed invalid mysql connection string", func(t *testing.T) {
		ResetInstance()

		os.Setenv("MYSQL_CONNECTION_URL", ":@invalid")
		_, err := GetDB(test.DB_SCHEMA)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to open database")

		os.Unsetenv("MYSQL_CONNECTION_URL")
	})

	t.Run("failed ping", func(t *testing.T) {
		ResetInstance()

		os.Setenv("MYSQL_CONNECTION_URL", "root:12345@tcp(127.0.0.1:3306)/casino?parseTime=true")
		_, err := GetDB(test.DB_SCHEMA)

		require.Error(t, err)
		require.ErrorContains(t, err, "failed to ping database")

		os.Unsetenv("MYSQL_CONNECTION_URL")
	})

	t.Run("failed prepad statement ", func(t *testing.T) {
		ResetInstance()

		_, err := GetDB(test.WRONG_DB_SCHEMA)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to prepare insert transaction statement")
	})

	t.Run("succesful singleton pattern", func(t *testing.T) {
		ResetInstance()

		// First call should succeed
		db1, err := GetDB(test.DB_SCHEMA)
		require.NoError(t, err)
		require.NotNil(t, db1)

		// Second call should return the same instance
		db2, err := GetDB(test.DB_SCHEMA)
		require.NoError(t, err)
		require.NotNil(t, db2)

		// Verify it's the same instance
		assert.Equal(t, db1, db2)
	})
}

// TestInsertTransaction tests the InsertTransaction method
func TestInsertTransaction(t *testing.T) {

	db, _ := GetDB(test.DB_SCHEMA)
	defer db.Close()

	t.Run("successful insert with valid data", func(t *testing.T) {
		err := db.InsertTransaction(test.USER_ID, test.TRANSACTION_TYPE, test.AMOUNT, time.Now())
		require.NoError(t, err)
	})

	t.Run("failed insert with invalid transaction type", func(t *testing.T) {
		err := db.InsertTransaction(test.USER_ID, test.WRONG_TRANSACTION_TYPE, test.AMOUNT, time.Now())
		require.Error(t, err)
		require.ErrorContains(t, err, "Data truncated for column 'transaction_type'")
	})

}

// TestGetTransactions tests the GetTransactions method
func TestGetTransactions(t *testing.T) {

	userId := test.USER_ID
	transactionType := test.TRANSACTION_TYPE

	t.Run("successful get transaction", func(t *testing.T) {
		db, _ := GetDB(test.DB_SCHEMA)
		defer db.Close()

		rows, err := db.GetTransactions(t.Context(), &userId, &transactionType, 1)
		require.NoError(t, err)
		defer rows.Close()

		assert.Equal(t, rows.Next(), true)
	})

	t.Run("failed get transaction", func(t *testing.T) {
		db, _ := GetDB(test.DB_SCHEMA)
		db.Close()

		_, err := db.GetTransactions(t.Context(), &userId, &transactionType, 1)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to query transactions: sql: statement is closed")
	})

}

func TestClose(t *testing.T) {

	t.Run("successfully closed", func(t *testing.T) {
		db, _ := GetDB(test.DB_SCHEMA)
		err := db.Close()

		require.Nil(t, instance)
		require.Nil(t, err)
	})

}
