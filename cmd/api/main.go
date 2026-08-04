package main

import (
	"banking-api/internal/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/health", handlers.HandleHealth)

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
