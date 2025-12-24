package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionRequest struct {
	WalletID      uuid.UUID       `json:"walletId"`
	OperationType string          `json:"operationType"`
	Amount        decimal.Decimal `json:"amount"`
}

type Wallet struct {
	ID      uuid.UUID       `json:"id"`
	Balance decimal.Decimal `json:"balance"`
}

/*type Transaction struct {
	ID            int       `json:"id"`
	WalletID      uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"createdAt"`
}
*/
