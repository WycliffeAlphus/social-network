package main

import (
	"fmt"
	"net/http"
	"backend/internal/routes"
)

func main() {
	routes.RegisterRoutes()
	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
