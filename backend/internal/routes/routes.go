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

	// Initialize Notification-related dependencies
	notificationRepo := repository.NewNotificationRepository(db)
	notificationService := service.NewNotificationService(notificationRepo, userRepo, groupRepo)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	groupHandler := &handler.GroupHandler{Service: groupService, NotificationService: notificationService}

	// Initialize Follower-related dependencies
	followerHandler := handler.NewFollowerHandler(db, notificationService)

	// Public routes (no authentication required)
	http.HandleFunc("/api/register", userHandler.Register)
	http.HandleFunc("/api/login", handler.LoginHandler)
	http.HandleFunc("/api/logout", handler.LogoutHandler)

	// http.Handle("/api/profile/", middlewares.AuthMiddleware(db, userHandler.Profile))
	http.Handle("/api/profile/", middlewares.AuthMiddleware(db, handler.ProfileHandler(db)))
	http.HandleFunc("/api/users/available", middlewares.AuthMiddleware(db, handler.GetFollowSuggestions(db)))
	http.HandleFunc("/api/users/follow", middlewares.AuthMiddleware(db, followerHandler.FollowUser))
	http.HandleFunc("/api/follow/accept", middlewares.AuthMiddleware(db, followerHandler.AcceptFollowRequest))
	http.HandleFunc("/api/follow/decline", middlewares.AuthMiddleware(db, followerHandler.DeclineFollowRequest))
	http.HandleFunc("/api/follow/cancel", middlewares.AuthMiddleware(db, followerHandler.CancelFollowRequest))
	http.HandleFunc("/api/follow-status/", middlewares.AuthMiddleware(db, followerHandler.GetFollowStatus()))
	http.HandleFunc("/api/followers/", middlewares.AuthMiddleware(db, followerHandler.GetFollowers()))
	http.HandleFunc("/api/following/", middlewares.AuthMiddleware(db, followerHandler.GetFollowing()))

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
		// Handle /api/groups/:id/invite endpoint
		if strings.Contains(r.URL.Path, "/invite") {
			middlewares.AuthMiddleware(db, http.HandlerFunc(groupHandler.InviteUserToGroup)).ServeHTTP(w, r)
			return
		}
		// Handle /api/groups/:id/events endpoint
		if strings.Contains(r.URL.Path, "/events") {
			middlewares.AuthMiddleware(db, http.HandlerFunc(groupHandler.CreateEvent)).ServeHTTP(w, r)
			return
		}
		http.Error(w, "Not found", http.StatusNotFound)
	})

	http.HandleFunc("/api/follow-requests", middlewares.AuthMiddleware(db, followerHandler.GetFollowRequests()))
	http.HandleFunc("/api/profile/update", middlewares.AuthMiddleware(db, handler.UpdateProfileHandler(db)))
	http.HandleFunc("/api/createpost", middlewares.AuthMiddleware(db, handler.CreatePost(db, notificationService)))

	// Comment routes
	http.HandleFunc("/api/posts/", middlewares.AuthMiddleware(db, handler.CommentHandler(db, notificationService)))
	http.HandleFunc("/api/feeds", middlewares.AuthMiddleware(db, handler.DashboardHandler(db)))
	http.HandleFunc("/api/reaction", middlewares.AuthMiddleware(db, handler.HandleReaction(db, notificationService)))

	// Notification routes
	http.HandleFunc("/api/notifications/read", middlewares.AuthMiddleware(db, notificationHandler.MarkNotificationsAsRead))
	http.HandleFunc("/api/notifications/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && strings.Contains(r.URL.Path, "/read") {
			middlewares.AuthMiddleware(db, notificationHandler.MarkNotificationAsRead).ServeHTTP(w, r)
		} else if r.Method == http.MethodGet {
			middlewares.AuthMiddleware(db, notificationHandler.GetNotifications).ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/api/notifications", middlewares.AuthMiddleware(db, notificationHandler.GetNotifications))
}
