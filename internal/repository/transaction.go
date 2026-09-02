package repository

import (
	"context"
	"fmt"
	"time"

	"banking-api/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	query := `
	INSERT INTO transactions (account_id, type, amount, balance_after, description, reference, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING transaction_id 
	`
	now := time.Now()

	var req models.Transaction
	err := r.db.QueryRow(
		ctx,
		query,
		req.AccountID,
		req.Type,
		req.Amount,
		req.BalanceAfter,
		req.Description,
		req.Reference,
		req.Reference,
		now,
	).Scan(&req.TransactionID)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	transaction.CreatedAt = now
	return nil
}
