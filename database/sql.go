package database

import (
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
	conn *sql.DB
}

var (
	instance *Database
	once     sync.Once
)

// GetDB returns a singleton MySQL database connection
func GetDB() (*Database, error) {
	var e error
	once.Do(func() {

		// Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		// Load mysql connection stirng from .env
		connectionString := os.Getenv("MYSQL_CONNECTION_URL")
		conn, err := sql.Open("mysql", connectionString)
		if err != nil {
			e = fmt.Errorf("failed to open database: %w", err)
			return
		}

		// Test the connection
		err = conn.Ping()
		if err != nil {
			e = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		// Configure connection pool settings
		conn.SetMaxOpenConns(25)
		conn.SetMaxIdleConns(25)
		conn.SetConnMaxLifetime(5 * time.Minute)

		instance = &Database{
			conn: conn,
		}
	})

	return instance, e
}

// Close the database connection
func (db *Database) Close() {
	if err := db.conn.Close(); err != nil {
		log.Printf("failed to close database")
	}
}
