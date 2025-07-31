package routes

import (
	"database/sql"
	"net/http"

	"backend/internal/handler"
	"backend/internal/middlewares"
	"backend/internal/repository"
	"backend/internal/service"
)

// RegisterRoutes sets up the HTTP routes for the API endpoints.
func RegisterRoutes(db *sql.DB) {
	// Initialize User-related dependencies
	userRepo := &repository.UserRepository{DB: db}
	userService := &service.UserService{Repo: userRepo}
	userHandler := &handler.UserHandler{Service: userService}

	groupRepo := repository.NewGroupRepository(db)
	groupService := service.NewGroupService(groupRepo)
	groupHandler := &handler.GroupHandler{Service: groupService}

	// Public routes (no authentication required)
	http.HandleFunc("/api/register", userHandler.Register)
	http.HandleFunc("/api/login", handler.LoginHandler)
	http.HandleFunc("/api/logout", handler.LogoutHandler)

	// http.Handle("/api/profile/", middlewares.AuthMiddleware(db, userHandler.Profile))
	http.Handle("/api/profile/", middlewares.AuthMiddleware(db, handler.ProfileHandler(db)))
	http.Handle("/api/dashboard", middlewares.AuthMiddleware(db, handler.DashboardHandler))
	http.HandleFunc("/api/users/available", middlewares.AuthMiddleware(db, handler.GetFollowSuggestions(db)))
	http.HandleFunc("/api/users/follow", middlewares.AuthMiddleware(db, handler.FollowUser(db)))
	http.HandleFunc("/api/follow/accept", middlewares.AuthMiddleware(db, handler.AcceptFollowRequest(db)))
	http.HandleFunc("/api/follow/decline", middlewares.AuthMiddleware(db, handler.DeclineFollowRequest(db)))
	http.HandleFunc("/api/followers/", middlewares.AuthMiddleware(db, handler.GetFollowers(db)))
	http.HandleFunc("/api/following/", middlewares.AuthMiddleware(db, handler.GetFollowing(db)))

	// This uses the initialized groupHandler instance and its CreateGroup method.
	http.HandleFunc("/api/groups", middlewares.AuthMiddleware(db, groupHandler.CreateGroup))
	http.HandleFunc("/api/follow-requests", middlewares.AuthMiddleware(db, handler.GetFollowRequests(db)))
}
