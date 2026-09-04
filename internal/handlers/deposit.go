package handlers

import (
	"banking-api/internal/middleware"
	"banking-api/internal/models"
	"banking-api/internal/repository"
	"banking-api/internal/validator"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func HandleDeposit(accountRepo *repository.AccountRepository, transactionRepo *repository.TransactionRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// POST Method
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"error":"method not allowed"}`))
			return
		}

		// Get userID from Context
		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}

		// Get accountID from URL Query
		accountIDStr := r.URL.Query().Get("account_id")
		accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"could not get account id"}`))
			return
		}

		// Bind data from request and parse to req struct
		var req models.DepositRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid request format"}`))
			return
		}

		// Validate the data
		if err := validator.Validate.Struct(req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"failed data validation"}`))
			return
		}

		// Get account from database by it's ID
		account, err := accountRepo.GetByID(r.Context(), accountID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to get account id"}`))
			return
		}

		// Ownership check
		if account.UserID != userID {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error":"forbidden"}`))
			return
		}

		// Calculate nee balance after deposit
		newBalance := account.Balance + req.Amount

		// Update the account balance after deposit
		err = accountRepo.UpdateBalance(r.Context(), accountID, newBalance)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to update balance"}`))
			return
		}

		// Create instance for transaction
		transaction := models.Transaction{
			AccountID:    accountID,
			Type:         "deposit",
			Amount:       req.Amount,
			BalanceAfter: newBalance,
			Description:  req.Description,
			Reference:    fmt.Sprintf("DEP-%d-%d", accountID, time.Now().UnixNano()),
		}

		// Create transaction
		err = transactionRepo.Create(r.Context(), &transaction)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to process deposit"}`))
			return
		}

		// Update account balance in struct
		account.Balance = newBalance

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":     "deposit successful",
			"account":     account,
			"transaction": transaction,
		})
	}
}
