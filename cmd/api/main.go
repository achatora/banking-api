package main

import (
	"fmt"
	"net/http"

	"banking-api/internal/handlers"
	"banking-api/internal/middleware"
	"banking-api/internal/router"
)

func main() {

	r := router.New()
	r.Handle("/", middleware.Logging(handlers.HandleHome))
	r.Handle("/health", middleware.Logging(handlers.HandleHealth))

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
