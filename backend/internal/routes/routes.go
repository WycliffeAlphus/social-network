package routes

import (
	"backend/internal/handler"
	"backend/internal/middlewares"
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

	// Public routes (no authentication required)
	http.HandleFunc("/api/register", userHandler.Register)
	http.HandleFunc("/api/login", handler.LoginHandler)
	http.HandleFunc("/api/logout", handler.LogoutHandler)

	// Protected routes (authentication required)
	http.Handle("/api/profile/", middlewares.AuthMiddleware(db, handler.ProfileHandler(db)))
	http.Handle("/api/dashboard", middlewares.AuthMiddleware(db, handler.DashboardHandler))
	http.HandleFunc("/api/users/available", middlewares.AuthMiddleware(db, handler.GetFollowSuggestions(db)))
	http.HandleFunc("/api/users/follow", middlewares.AuthMiddleware(db, handler.FollowUser(db)))
	http.HandleFunc("/api/follow/accept", middlewares.AuthMiddleware(db, handler.AcceptFollowRequest(db)))
	http.HandleFunc("/api/follow/decline", middlewares.AuthMiddleware(db, handler.DeclineFollowRequest(db)))
	http.HandleFunc("/api/follow/cancel", middlewares.AuthMiddleware(db, handler.CancelFollowRequest(db)))
	http.HandleFunc("/api/follow-status/", middlewares.AuthMiddleware(db, handler.GetFollowStatus(db)))
	http.HandleFunc("/api/followers/", middlewares.AuthMiddleware(db, handler.GetFollowers(db)))
	http.HandleFunc("/api/following/", middlewares.AuthMiddleware(db, handler.GetFollowing(db)))
	http.HandleFunc("/api/follow-requests", middlewares.AuthMiddleware(db, handler.GetFollowRequests(db)))
	http.HandleFunc("/api/profile/update", middlewares.AuthMiddleware(db, handler.UpdateProfileHandler(db)))
	http.HandleFunc("/api/createpost", middlewares.AuthMiddleware(db, handler.CreatePost(db)))

	// Comment routes
	http.HandleFunc("/api/posts/", middlewares.AuthMiddleware(db, handler.CommentHandler(db)))
}
