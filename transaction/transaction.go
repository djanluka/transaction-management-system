package transaction

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	BET = "bet"
	WIN = "win"
)

var TransactionTypes = []string{
	BET,
	WIN,
}

type Transaction struct {
	UserId          int       `json:"user_id"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	Timestamp       time.Time `json:"timestamp"`
}

func NewTransaction() Transaction {
	return Transaction{
		UserId:          setUserId(),
		TransactionType: setTransactionType(),
		Amount:          setAmount(),
		Timestamp:       setTimestamp(),
	}
}

func setUserId() int {
	return rand.Intn(5) + 1
}

func setTransactionType() string {
	return TransactionTypes[rand.Intn(len(TransactionTypes))]
}

func setAmount() float64 {
	return math.Round(rand.Float64()*100) / 100
}

func setTimestamp() time.Time {
	return time.Now()
}

func (t Transaction) String() string {
	return fmt.Sprintf("{user_id: %d, transaction_type: %s, amount: %.2f, timestamp: %s}", t.UserId, t.TransactionType, t.Amount, t.Timestamp.Format(time.RFC1123))
}
