package routes

import (
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"
	"database/sql"
	"net/http"
)

// RegisterRoutes sets up all the user-related web routes for our app
// It creates the different parts that handle user requests step by step
func RegisterRoutes(db *sql.DB) {
	// Step 1: Create the part that talks to the database
	userRepo := &repository.UserRepository{DB: db}

	// Step 2: Create the part that handles the business rules and checks
	userService := &service.UserService{Repo: userRepo}

	// Step 3: Create the part that handles web requests from users
	userHandler := &handler.UserHandler{Service: userService}

	// Step 4: Tell our app which web address should do what
	// When someone visits /api/register, run the Register function
	http.HandleFunc("/api/register", userHandler.Register)
}
