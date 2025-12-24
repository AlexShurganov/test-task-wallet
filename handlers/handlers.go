package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"wallet-service/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Job struct {
	req    models.TransactionRequest
	respCh chan JobResponse
}

type JobResponse struct {
	NewBalance decimal.Decimal
	Err        error
}

var (
	db   *sql.DB
	jobs chan Job
)

func SetDB(database *sql.DB) {
	db = database
}

func InitWorkerPool(numWorkers int) {
	jobs = make(chan Job, 1000)

	for i := 0; i < numWorkers; i++ {
		go worker()
	}
}

func worker() {
	for job := range jobs {
		newBalance, err := processTransaction(job.req)
		job.respCh <- JobResponse{NewBalance: newBalance, Err: err}
	}
}

func processTransaction(req models.TransactionRequest) (decimal.Decimal, error) {
	changeAmount := req.Amount
	if req.OperationType == "WITHDRAW" {
		changeAmount = req.Amount.Neg()
	}

	var newBalance decimal.Decimal

	query := `
		UPDATE wallets 
		SET balance = balance + $2 
		WHERE id = $1 AND (balance + $2 >= 0)
		RETURNING balance`

	err := db.QueryRow(query, req.WalletID, changeAmount).Scan(&newBalance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return decimal.Zero, errors.New("wallet not found or insufficient funds")
		}
		return decimal.Zero, err
	}

	_, _ = db.Exec(
		"INSERT INTO transactions (wallet_id, operation_type, amount) VALUES ($1, $2, $3)",
		req.WalletID, req.OperationType, req.Amount,
	)

	return newBalance, nil
}

func Transaction(c *gin.Context) {
	var req models.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if req.OperationType != "DEPOSIT" && req.OperationType != "WITHDRAW" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation type"})
		return
	}
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}

	respCh := make(chan JobResponse)

	job := Job{
		req:    req,
		respCh: respCh,
	}

	jobs <- job

	response := <-respCh

	if response.Err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Transaction successful",
		"walletId":   req.WalletID,
		"newBalance": response.NewBalance,
	})
}

func WalletBalance(c *gin.Context) {
	walletIDStr := c.Param("id")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var wallet models.Wallet
	err = db.QueryRow("SELECT id, balance FROM wallets WHERE id = $1", walletID).Scan(&wallet.ID, &wallet.Balance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, wallet)
}
