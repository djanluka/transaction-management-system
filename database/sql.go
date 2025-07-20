package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"transaction-management-system/transaction"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// DB is a singleton struct that holds the database connection and prepared statements.
type Database struct {
	conn                      *sql.DB
	insertTransactionPrepStmt *sql.Stmt
}

var (
	instance *Database
	once     sync.Once
)

// GetDB returns a singleton MySQL database connection
func GetDB() (*Database, error) {
	var initError error

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load mysql connection stirng from .env
	connectionString := os.Getenv("MYSQL_CONNECTION_URL")

	once.Do(func() {

		conn, err := sql.Open("mysql", connectionString)
		if err != nil {
			initError = fmt.Errorf("failed to open database: %w", err)
			return
		}

		// Test the connection
		err = conn.Ping()
		if err != nil {
			initError = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		// Configure connection pool settings
		conn.SetMaxOpenConns(25)
		conn.SetMaxIdleConns(25)
		conn.SetConnMaxLifetime(5 * time.Minute)

		// Prepare the transaction insert statement
		insertTransactionPrepStmt, err := conn.Prepare("INSERT INTO casino.transactions (user_id, transaction_type, amount, timestamp) VALUES (?, ?, ?, ?)")
		if err != nil {
			initError = fmt.Errorf("failed to prepare insert transaction statement: %w", err)
			return
		}

		instance = &Database{
			conn:                      conn,
			insertTransactionPrepStmt: insertTransactionPrepStmt,
		}
	})

	return instance, initError
}

// InsertTransaction inserts a new transaction record
func (db *Database) InsertTransaction(tr *transaction.Transaction) error {
	_, err := db.insertTransactionPrepStmt.Exec(
		tr.UserId,
		tr.TransactionType,
		tr.Amount,
		tr.Timestamp,
	)
	return err
}

// Close the database connection
func (db *Database) Close() {
	if err := db.insertTransactionPrepStmt.Close(); err != nil {
		log.Printf("failed to close insert prepared statement")
	}

	if err := db.conn.Close(); err != nil {
		log.Printf("failed to close database")
	}
}
