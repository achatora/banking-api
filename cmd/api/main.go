package main

import (
	"context"
	"fmt"
	"net/http"

	"banking-api/internal/config"
	"banking-api/internal/database"
	"banking-api/internal/handlers"
	"banking-api/internal/middleware"
	"banking-api/internal/router"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

	cfg := config.Load()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := database.New(connString)
	if err != nil {
		fmt.Println("Failed to initialise database:", err)
		return
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		fmt.Println("Failed to ping database:", err)
		return
	}

	r := router.New()
	r.Handle("/", middleware.Logging(handlers.HandleHome))
	r.Handle("/health", middleware.Logging(handlers.HandleHealth))

	fmt.Printf("Server starting on: %s\n", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
