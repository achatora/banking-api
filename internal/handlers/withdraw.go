package handlers

import (
	"banking-api/internal/middleware"
	"banking-api/internal/models"
	"banking-api/internal/repository"
	"banking-api/internal/validator"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HandleWithdraw(accountRepo *repository.AccountRepository, transactionRepo *repository.TransactionRepository) http.HandlerFunc {

	// POST Method
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"error":"method not allowed"}`))
			return
		}

		// Get userID from middleware Context
		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"no authorized"}`))
			return

		}

		// Get account ID from URL Query
		accountIDStr := r.URL.Query().Get("account_id")
		accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"could not get account id"}`))
			return
		}

		// Bind data and parse to req struct
		var req models.WithdrawRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid request format"}`))
			return
		}

		// Validate data
		if err := validator.Validate.Struct(req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"failed data validation"}`))
			return
		}

		// Get account from database by user's ID
		account, err := accountRepo.GetByID(r.Context(), accountID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"could not get account"}`))
			return
		}

		// Ownership check
		if account.UserID != userID {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error":"forbidden"}`))
			return
		}

		// Balance check
		if account.Balance < req.Amount {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"low account balance, cannot withdraw"}`))
			return
		}

		// Withdraw ammount from balance
		newBalance := account.Balance - req.Amount

		// Update balance in database
		err = accountRepo.UpdateBalance(r.Context(), accountID, newBalance)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to update account balance"}`))
			return
		}

		transaction := models.Transaction{
			AccountID:    accountID,
			Type:         "withdrawal",
			Amount:       req.Amount,
			BalanceAfter: newBalance,
			Description:  req.Description,
			Reference:    fmt.Sprintf("WTH-%d-%d", userID, time.Now().UnixNano()),
		}

		// Create transaction record in database (transaction history)
		err = transactionRepo.Create(r.Context(), &transaction)
		if err != nil {
			log.Printf("transaction record failed: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to create transaction record"}`))
			return
		}

		// Update account balance in struct
		account.Balance = newBalance

		// Return data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":     "withdrawal successful",
			"account":     account,
			"transaction": transaction,
		})
	}
}
