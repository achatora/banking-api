package models

import "time"

type Account struct {
	AccountID     int64     `json:"account_id"`
	UserID        int64     `json:"user_id"`
	AccountNumber string    `json:"account_number"`
	AccountType   string    `json:"account_type"`
	Balance       int64     `json:"balance"`
	Status        string    `json:"Status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateAccountRequest struct {
	AccountType string `json:"account_type" validate:"required,oneof=checking savings"`
	Balance     int64  `json:"balance" validate:"gte=0"`
}
