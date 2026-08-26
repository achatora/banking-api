package repository

import (
	"context"
	"fmt"
	"time"

	"banking-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
	query := `
	INSERT INTO accounts (user_id, account_number, account_type, balance, status, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING account_id
	`

	now := time.Now()
	err := r.db.QueryRow(
		ctx,
		query,
		account.UserID,
		account.AccountNumber,
		account.AccountType,
		account.Balance,
		account.Status,
		now,
		now,
	).Scan(&account.AccountID)

	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	account.CreatedAt = now
	account.UpdatedAt = now

	return nil
}

func (r *AccountRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Account, error) {
	query := `
	SELECT account_id, user_id, account_number, account_type, balance, status, created_at, updated_at
	FROM accounts
	WHERE user_id = $1
	ORDER BY created_at DESC 
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.AccountID,
			&account.UserID,
			&account.AccountNumber,
			&account.AccountType,
			&account.Balance,
			&account.Status,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan accounts: %w", err)
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate accounts: %w", err)
	}

	return accounts, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, accountID int64) (*models.Account, error) {
	query := `
	SELECT account_id, user_id, account_number, account_type, balance, status, created_at, updated_at
	FROM accounts
	WHERE account_id = $1
	`
	var account models.Account
	err := r.db.QueryRow(
		ctx,
		query,
		account.AccountID,
	).Scan(
		&account.AccountID,
		&account.UserID,
		&account.AccountNumber,
		&account.AccountType,
		&account.Balance,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get account id: %w", err)
	}

	return &account, nil
}
