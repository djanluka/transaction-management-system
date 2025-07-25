package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// DB is a singleton struct that holds the database connection and prepared statements.
type Database struct {
	conn                      *sql.DB
	insertTransactionPrepStmt *sql.Stmt
	getTransactionsPrepStmt   *sql.Stmt
}

var (
	instance *Database
	once     sync.Once
)

// GetDB returns a singleton MySQL database connection
func GetDB(schema string) (*Database, error) {
	var initError error

	once.Do(func() {
		envPath := os.Getenv("ENV_PATH")
		if envPath == "" {
			envPath = ".env" // default
		}

		// Load .env file
		err := godotenv.Load(envPath)
		if err != nil {
			initError = fmt.Errorf("failed to load .env file")
			return
		}
		// Load mysql connection stirng from .env
		connectionString := os.Getenv("MYSQL_CONNECTION_URL")

		// Connect to MySQL database
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
		insertTransactionPrepStmt, err := conn.Prepare(fmt.Sprintf("INSERT INTO %s.transactions (user_id, transaction_type, amount, timestamp) VALUES (?, ?, ?, ?)", schema))
		if err != nil {
			initError = fmt.Errorf("failed to prepare insert transaction statement: %w", err)
			return
		}
		// Prepate get transactions statement
		getTransactionsPrepStmt, err := conn.Prepare(fmt.Sprintf(`
			SELECT user_id, transaction_type, amount, timestamp 
			FROM %s.transactions 
			WHERE (? IS NULL OR user_id = ?)
			AND (? IS NULL OR transaction_type = ?)
			ORDER BY timestamp DESC
			LIMIT ?
		`, schema))
		if err != nil {
			initError = fmt.Errorf("failed to prepare get transactions statement: %w", err)
			return
		}

		instance = &Database{
			conn:                      conn,
			insertTransactionPrepStmt: insertTransactionPrepStmt,
			getTransactionsPrepStmt:   getTransactionsPrepStmt,
		}
	})

	return instance, initError
}

// InsertTransaction inserts a new transaction record
func (db *Database) InsertTransaction(userId int, transactionType string, amount float64, timestamp time.Time) error {
	_, err := db.insertTransactionPrepStmt.Exec(
		userId,
		transactionType,
		amount,
		timestamp,
	)
	return err
}

func (db *Database) GetTransactions(ctx context.Context, userId *int, transactionType *string, limit int) (*sql.Rows, error) {
	var userIdVal, typeVal interface{}
	if userId != nil {
		userIdVal = *userId
	}
	if transactionType != nil {
		typeVal = *transactionType
	}

	rows, err := db.getTransactionsPrepStmt.QueryContext(ctx, userIdVal, userIdVal, typeVal, typeVal, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	return rows, nil
}

// Close the database connection
func (db *Database) Close() error {
	if err := db.insertTransactionPrepStmt.Close(); err != nil {
		log.Printf("Failed to close insert prepared statement")
		return err
	}

	if err := db.getTransactionsPrepStmt.Close(); err != nil {
		log.Printf("Failed to close insert prepared statement")
		return err
	}

	if err := db.conn.Close(); err != nil {
		log.Printf("Failed to close database")
		return err
	}

	ResetInstance()
	log.Println("Database closed succesfully")
	return nil
}

func ResetInstance() {
	instance = nil
	once = sync.Once{}
}
