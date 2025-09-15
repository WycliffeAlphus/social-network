package routes

import (
	"database/sql"
	"net/http"
	"strings"

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
	http.HandleFunc("/api/users/available", middlewares.AuthMiddleware(db, handler.GetFollowSuggestions(db)))
	http.HandleFunc("/api/users/follow", middlewares.AuthMiddleware(db, handler.FollowUser(db)))
	http.HandleFunc("/api/follow/accept", middlewares.AuthMiddleware(db, handler.AcceptFollowRequest(db)))
	http.HandleFunc("/api/follow/decline", middlewares.AuthMiddleware(db, handler.DeclineFollowRequest(db)))
	http.HandleFunc("/api/follow/cancel", middlewares.AuthMiddleware(db, handler.CancelFollowRequest(db)))
	http.HandleFunc("/api/follow-status/", middlewares.AuthMiddleware(db, handler.GetFollowStatus(db)))
	http.HandleFunc("/api/followers/", middlewares.AuthMiddleware(db, handler.GetFollowers(db)))
	http.HandleFunc("/api/following/", middlewares.AuthMiddleware(db, handler.GetFollowing(db)))

	groupsHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middlewares.AuthMiddleware(db, http.HandlerFunc(groupHandler.GetGroups)).ServeHTTP(w, r)
		case http.MethodPost:
			middlewares.AuthMiddleware(db, http.HandlerFunc(groupHandler.CreateGroup)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}

	http.HandleFunc("/api/groups", groupsHandler)

	// Group join request endpoints
	http.HandleFunc("/api/groups/", func(w http.ResponseWriter, r *http.Request) {
		// Handle /api/groups/:id/join endpoint
		if strings.Contains(r.URL.Path, "/join") {
			if r.URL.Query().Get("action") == "accept" {
				middlewares.AuthMiddleware(db, http.HandlerFunc(groupHandler.AcceptJoinRequest)).ServeHTTP(w, r)
			} else {
				middlewares.AuthMiddleware(db, http.HandlerFunc(groupHandler.JoinGroupRequest)).ServeHTTP(w, r)
			}
			return
		}

		// Handle /api/groups/:id/events endpoint
		if strings.Contains(r.URL.Path, "/events") {
			middlewares.AuthMiddleware(db, http.HandlerFunc(groupHandler.CreateGroupEvent)).ServeHTTP(w, r)
			return
		}

		// Handle /api/groups/:id/view endpoint
		if strings.Contains(r.URL.Path, "/view") {
			middlewares.AuthMiddleware(db, http.HandlerFunc(groupHandler.ViewGroup)).ServeHTTP(w, r)
			return
		}

		http.Error(w, "Not found", http.StatusNotFound)
	})

	http.HandleFunc("/api/follow-requests", middlewares.AuthMiddleware(db, handler.GetFollowRequests(db)))
	http.HandleFunc("/api/profile/update", middlewares.AuthMiddleware(db, handler.UpdateProfileHandler(db)))
	http.HandleFunc("/api/createpost", middlewares.AuthMiddleware(db, handler.CreatePost(db)))

	// Comment routes
	http.HandleFunc("/api/posts/", middlewares.AuthMiddleware(db, handler.CommentHandler(db)))
	http.HandleFunc("/api/feeds", middlewares.AuthMiddleware(db, handler.DashboardHandler(db)))
	http.HandleFunc("/api/reaction", middlewares.AuthMiddleware(db, handler.HandleReaction(db)))

}
