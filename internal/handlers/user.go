package handlers

import (
	"banking-api/internal/auth"
	"banking-api/internal/config"
	"banking-api/internal/middleware"
	"banking-api/internal/models"
	"banking-api/internal/password"
	"banking-api/internal/repository"
	"banking-api/internal/validator"
	"encoding/json"
	"net/http"
	"strconv"
)

func HandleRegister(userRepo *repository.UserRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"error":"method not allowed"}`))
			return
		}

		// Bind our data from request body
		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid request"}`))
			return
		}

		if err := validator.Validate.Struct(req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"failed data validation"}`))
			return
		}

		// Check if email already exists
		existing, err := userRepo.GetByEmail(r.Context(), req.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to check email"}`))
			return
		}

		if existing != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"error":"email already exists"}`))
			return
		}

		// password hashing
		hashed, err := password.Hash(req.Password)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to hash password"}`))
			return
		}

		// Sruct with data from context
		user := &models.User{
			Email:        req.Email,
			PasswordHash: hashed,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
		}

		// Handle erorr if not nil when creating user
		if err := userRepo.Create(r.Context(), user); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to create user"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success":"user account created"}`))
	}
}

func HandleLogin(userRepo *repository.UserRepository, cfg *config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Mehtod Post
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"error":"method not allowed"}`))
			return
		}

		// Bind data and decode into req struct
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid request"}`))
			return
		}

		// Validate the data
		if err := validator.Validate.Struct(req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"failed data validation"}`))
			return
		}

		// Check if account exists and handle error
		user, err := userRepo.GetByEmail(r.Context(), req.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to check email"}`))
			return
		}

		if user == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"invalid credentials"}`))
			return
		}

		// Comapare login password to hash
		if err := password.Compare(user.PasswordHash, req.Password); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"invalid credentials"}`))
			return
		}

		// Create jwt token
		token, err := auth.GenerateJWT(user.UserID, user.Email, cfg.JWTSecret)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to generate token"}`))
			return
		}

		// return token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}

func HandleProfile(userRepo *repository.UserRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Get user ID from context (set by auth middleware)
		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}

		// Get user from database
		user, err := userRepo.GetByID(r.Context(), userID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"failed to get user"}`))
			return
		}

		if user == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"user not found"}`))
			return
		}

		// Return user data (password hash is hidden by json: "-")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func HandleGetAccountByID(accountRepo *repository.AccountRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"error":"method not allowed"}`))
			return
		}

		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}

		idStr := r.URL.Query().Get("id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid request format"}`))
			return
		}

		account, err := accountRepo.GetByID(r.Context(), id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"could not get account"}`))
			return
		}

		if account == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"account not found"}`))
			return
		}

		if account.UserID != userID {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error":"forbidden"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(account)
	}
}
