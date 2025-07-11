package routes

import (
	"backend/internal/handler"
	"net/http"
)

// RegisterRoutes sets up the HTTP routes for the API endpoints.
func RegisterRoutes() {
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)
}
