package routes

import (
	"backend/internal/handler"
	"net/http"
)

// RegisterRoutes sets up the HTTP routes for the API endpoints.
func RegisterRoutes() {
	http.HandleFunc("/api/login", handler.LoginHandler)
	http.HandleFunc("/api/logout", handler.LogoutHandler)
}
