package routes

import (
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"
	"database/sql"
	"net/http"
)

// RegisterRoutes sets up the HTTP routes for the API endpoints.
func RegisterRoutes(db *sql.DB) {
	userRepo := &repository.UserRepository{DB: db}

	// Step 2: Create the part that handles the business rules and checks
	userService := &service.UserService{Repo: userRepo}

	// Step 3: Create the part that handles web requests from users
	userHandler := &handler.UserHandler{Service: userService}
	http.HandleFunc("/api/register", userHandler.Register)
	http.HandleFunc("/api/login", handler.LoginHandler)
	http.HandleFunc("/api/logout", handler.LogoutHandler)
}
