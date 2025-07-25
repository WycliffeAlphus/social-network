package routes

import (
	"backend/internal/handler"
	"backend/internal/middlewares"
	"backend/internal/repository"
	"backend/internal/service"
	"database/sql"
	"net/http"
	"strings"
)

// RegisterRoutes sets up the HTTP routes for the API endpoints.
func RegisterRoutes(db *sql.DB) {
	userRepo := &repository.UserRepository{DB: db}

	// Step 2: Create the part that handles the business rules and checks
	userService := &service.UserService{Repo: userRepo}

	// Step 3: Create the part that handles web requests from users
	userHandler := &handler.UserHandler{Service: userService}

	// Public routes (no authentication required)
	http.HandleFunc("/api/register", userHandler.Register)
	http.HandleFunc("/api/login", handler.LoginHandler)
	http.HandleFunc("/api/logout", handler.LogoutHandler)

	// Protected routes (authentication required)
	http.Handle("/api/profile", middlewares.AuthMiddleware(db, handler.ProfileHandler))
	http.Handle("/api/profile/update", middlewares.AuthMiddleware(db, handler.UpdateProfileHandler))
	http.Handle("/api/dashboard", middlewares.AuthMiddleware(db, handler.DashboardHandler))
	// Follower routes - only GET endpoints for followers and following lists
	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Routes that require authentication (for viewing lists)
		if strings.HasSuffix(path, "/followers") && r.Method == http.MethodGet {
			middlewares.AuthMiddleware(db, handler.GetFollowers(db))(w, r)
			return
		}

		if strings.HasSuffix(path, "/following") && r.Method == http.MethodGet {
			middlewares.AuthMiddleware(db, handler.GetFollowing(db))(w, r)
			return
		}

		http.NotFound(w, r)
	})
}
