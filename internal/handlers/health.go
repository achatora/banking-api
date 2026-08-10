package handlers

import (
	"banking-api/internal/database"
	"net/http"
)

func HandleHealth(db *database.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		err := db.Ping(r.Context())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"database unreachable"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
}
