package handlers

import (
	"banking-api/internal/middleware"
	"banking-api/internal/models"
	"banking-api/internal/repository"
	"banking-api/internal/validator"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func HandleAccounts(accountRepo *repository.AccountRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Check Method type
		if r.Method == http.MethodPost { // POST METHOD
			// Get userID from (middleware) Context
			userID, ok := middleware.GetUserID(r.Context())
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"unauthorized"}`))
				return
			}

			// Bind the data from Context
			var req models.CreateAccountRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"invalid request Body"}`))
				return
			}

			// Validate data
			if err := validator.Validate.Struct(req); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"failed data validation"}`))
				return
			}

			// Create unique bank account number
			accountNumber := fmt.Sprintf("ACC-%d%d", time.Now().Unix(), time.Now().Nanosecond()%1000)

			// Create account
			account := &models.Account{
				UserID:        userID,
				AccountNumber: accountNumber,
				AccountType:   req.AccountType,
				Balance:       req.Balance,
				Status:        "active",
			}

			// Save account to database
			err = accountRepo.Create(r.Context(), account)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"failed to create account"}`))
				return
			}

			// Return the account
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(account)

		} else if r.Method == http.MethodGet { // GET METHOD

			// Get ID from (middleware) Context
			userID, ok := middleware.GetUserID(r.Context())
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"unauthorized"}`))
				return
			}

			// Get user's bank accounts by userID
			accounts, err := accountRepo.GetByUserID(r.Context(), userID)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"could not retrieve accounts"}`))
				return
			}

			// If no accounts found return empty slice
			if accounts == nil {
				accounts = []models.Account{}
			}

			// Return accounts
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(accounts)

		} else { // If invalid Method

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		}
	}
}
