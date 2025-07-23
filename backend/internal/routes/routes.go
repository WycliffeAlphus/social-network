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

	// Initialize followers handler
	followersHandler := handler.NewFollowersHandler(db)

	// Protected routes (authentication required)
	http.Handle("/api/profile", middlewares.AuthMiddleware(db, handler.ProfileHandler))
	http.Handle("/api/profile/update", middlewares.AuthMiddleware(db, handler.UpdateProfileHandler))
	http.Handle("/api/dashboard", middlewares.AuthMiddleware(db, handler.DashboardHandler))
	http.HandleFunc("/api/users/available", middlewares.AuthMiddleware(db, handler.GetFollowSuggestions(db)))
	http.HandleFunc("/api/users/follow", middlewares.AuthMiddleware(db, handler.FollowUser(db)))

	// Follower routes - these need custom routing due to path parameters
	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Routes that require authentication
		if strings.HasSuffix(path, "/follow") && (r.Method == http.MethodPost || r.Method == http.MethodDelete) {
			middlewares.AuthMiddleware(db, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					followersHandler.FollowUser(w, r)
				} else if r.Method == http.MethodDelete {
					followersHandler.UnfollowUser(w, r)
				}
			}))(w, r)
			return
		}

		// Follow request management routes (require authentication)
		if strings.Contains(path, "/follow-requests/") {
			middlewares.AuthMiddleware(db, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasSuffix(path, "/accept") {
					followersHandler.AcceptFollowRequest(w, r)
				} else if strings.HasSuffix(path, "/reject") {
					followersHandler.RejectFollowRequest(w, r)
				} else {
					http.NotFound(w, r)
				}
			}))(w, r)
			return
		}

		if strings.HasSuffix(path, "/follow-requests") {
			middlewares.AuthMiddleware(db, http.HandlerFunc(followersHandler.GetPendingRequests))(w, r)
			return
		}

		// Routes with optional authentication (for viewing lists)
		if strings.HasSuffix(path, "/followers") && r.Method == http.MethodGet {
			middlewares.OptionalAuth(db, http.HandlerFunc(followersHandler.GetFollowers))(w, r)
			return
		}

		if strings.HasSuffix(path, "/following") && r.Method == http.MethodGet {
			middlewares.OptionalAuth(db, http.HandlerFunc(followersHandler.GetFollowing))(w, r)
			return
		}

		if strings.HasSuffix(path, "/followers/count") && r.Method == http.MethodGet {
			middlewares.OptionalAuth(db, http.HandlerFunc(followersHandler.GetFollowerCounts))(w, r)
			return
		}

		http.NotFound(w, r)
	})
}
