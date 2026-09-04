package models

import "time"

type Transaction struct {
	TransactionID int64     `json:"transaction_id"`
	AccountID     int64     `json:"account_id"`
	Type          string    `json:"type"`
	Amount        int64     `json:"amount"`
	BalanceAfter  int64     `json:"balance_after"`
	Description   string    `json:"description,omitempty"`
	Reference     string    `json:"reference"`
	CreatedAt     time.Time `json:"created_at"`
}

type DepositRequest struct {
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	Description string `json:"description,omitempty"`
}

type WithdrawRequest struct {
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	Description string `json:"description,omitempty"`
}
