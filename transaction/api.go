package transaction

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
	"transaction-management-system/database"
)

type TransactionApi struct {
	Database *database.Database
}

func NewTransactionApi() *TransactionApi {
	db, err := database.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to Database: %v", err)
	}

	return &TransactionApi{
		Database: db,
	}
}

// GetTransactions handles GET requests for transaction data
func (tapi *TransactionApi) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Get optional user_id filter
	var userId *int
	if uid := query.Get("user_id"); uid != "" {
		user_id, err := strconv.Atoi(uid)
		if err != nil {
			http.Error(w, "Invalid user id conversion", http.StatusBadRequest)
			return
		}
		userId = &user_id
	}

	// Get optional transaction_type filter ("bet", "win", or empty for all)
	var transactionType *string
	if tt := query.Get("transaction_type"); tt != "" {
		if tt == "bet" || tt == "win" {
			transactionType = &tt
		} else {
			http.Error(w, "Invalid transaction_type. Must be 'bet' or 'win'", http.StatusBadRequest)
			return
		}
	}

	// Get optional limit parameter
	limit := math.MaxInt
	if l := query.Get("limit"); l != "" {
		lInt, err := strconv.Atoi(l)
		if err != nil || lInt < 1 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		limit = lInt
	}

	// Get transactions from database
	rows, err := tapi.Database.GetTransactions(r.Context(), userId, transactionType, limit)
	if err != nil {
		http.Error(w, "Failed to retrieve transactions: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Scan all transactions
	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.UserId, &t.TransactionType, &t.Amount, &t.Timestamp); err != nil {
			http.Error(w, "Failed to scan transaction: "+err.Error(), http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, t)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error processing transaction data", http.StatusInternalServerError)
		return
	}

	// Return as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// RegisterRoutes sets up the HTTP routes for the Transaction API
func (tapi *TransactionApi) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/transactions", tapi.GetTransactions)
}

func (tapi *TransactionApi) ListenAndServe(wg *sync.WaitGroup) {
	defer wg.Done()

	// Set up router
	mux := http.NewServeMux()
	tapi.RegisterRoutes(mux)

	// Configure server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
