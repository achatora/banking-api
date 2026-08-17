package handlers

import (
	"banking-api/internal/models"
	"banking-api/internal/password"
	"banking-api/internal/repository"
	"banking-api/internal/validator"
	"encoding/json"
	"net/http"
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
