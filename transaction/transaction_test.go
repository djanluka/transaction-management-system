package transaction

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
	"transaction-management-system/database"

	"github.com/stretchr/testify/require"
)

func TestNewTransaction(t *testing.T) {
	t.Run("succesful new transaction", func(t *testing.T) {
		tr := NewTransaction()
		require.NotEmpty(t, tr)
	})
}

func TestGetUserId(t *testing.T) {
	t.Run("successful get user id", func(t *testing.T) {
		id := getUserId()
		require.Greater(t, id, 0)
	})
}

func TestGetTransactionType(t *testing.T) {
	t.Run("successful get transaction type", func(t *testing.T) {
		ttype := getTransactionType()
		isValidType := ttype == BET || ttype == WIN

		require.Equal(t, isValidType, true)
	})
}

func TestGetAmount(t *testing.T) {
	t.Run("successful get amount", func(t *testing.T) {
		a := getAmount()
		require.GreaterOrEqual(t, a, 0.0)
	})
}

func TestGetTimestamp(t *testing.T) {
	t.Run("successful get user id", func(t *testing.T) {
		before := time.Now()
		ts := getTimestamp()
		after := time.Now()

		require.True(t, ts.After(before) || ts.Equal(before),
			"timestamp should be after or equal to time before call")
		require.True(t, ts.Before(after) || ts.Equal(after),
			"timestamp should be before or equal to time after call")
	})
}

func TestNewTransactionApi(t *testing.T) {
	t.Run("failed new transaction api - db", func(t *testing.T) {
		os.Setenv("ENV_PATH", "")

		_, err := NewTransactionApi()
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to load .env file")

		os.Setenv("ENV_PATH", "../.env")
		database.ResetInstance()
	})
	t.Run("successful new transaction api", func(t *testing.T) {
		tapi, err := NewTransactionApi()

		require.Nil(t, err)
		require.NotNil(t, tapi)
	})
}

func TestGetTransactions(t *testing.T) {

	tapi, _ := NewTransactionApi()

	// Create a test HTTP server
	srv := httptest.NewServer(http.HandlerFunc(tapi.GetTransactions))
	defer srv.Close()

	t.Run("failed: invalid user id", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/transactions?user_id=abc")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		bodyStr := string(body)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Contains(t, bodyStr, "Invalid user id conversion")
	})

	t.Run("failed: invalid transaction type", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/transactions?transaction_type=abc")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		bodyStr := string(body)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Contains(t, bodyStr, "Invalid transaction_type")
	})

	t.Run("failed: invalid limit", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/transactions?limit=abc")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		bodyStr := string(body)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		require.Contains(t, bodyStr, "Invalid limit parameter")
	})

	t.Run("failed: db get transactions", func(t *testing.T) {
		// Close DB to produce an error
		tapi.Database.Close()
		defer database.ResetInstance()

		resp, err := http.Get(srv.URL + "/transactions?limit=1")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		bodyStr := string(body)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		require.Contains(t, bodyStr, "Failed to retrieve transactions")
	})

	tapi, _ = NewTransactionApi()
	srv = httptest.NewServer(http.HandlerFunc(tapi.GetTransactions))
	t.Run("succesful request", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/transactions?user_id=123&transaction_type=bet")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestRegisterRoutes(t *testing.T) {
	// Create mux for tapi
	tapi, _ := NewTransactionApi()
	mux := http.NewServeMux()
	tapi.RegisterRoutes(mux)

	// Verify the route was registered by making a test request
	req := httptest.NewRequest("GET", "/transactions", nil)
	rr := httptest.NewRecorder()

	// ServeMux will find and call the registered handler
	mux.ServeHTTP(rr, req)

	require.NotEqual(t, http.StatusNotFound, rr.Code,
		"Expected handler to be registered for /transactions")
}

func TestListenAndServe(t *testing.T) {
	var wg sync.WaitGroup
	tapi, _ := NewTransactionApi()

	wg.Add(1)
	go tapi.ListenAndServe(&wg)

	// Find the process and kill it sub-goroutine
	process, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)

	go process.Signal(syscall.SIGTERM)

	wg.Wait()
}
