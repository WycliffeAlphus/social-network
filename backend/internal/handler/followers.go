package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// FollowersHandler handles all follower-related operations
type FollowersHandler struct {
	Service *service.FollowerService
}

// GetFollowers handles GET /users/:id/followers
func (h *FollowersHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path
	userID := extractUserIDFromPath(r.URL.Path, "followers")
	if userID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Handle "me" by resolving to current user's ID
	if userID == "me" {
		if user, ok := context.GetUser(r.Context()); ok {
			userID = user.ID
		} else {
			http.Error(w, "Authentication required to view your own followers", http.StatusUnauthorized)
			return
		}
	}

	// Get requesting user ID from context (may be empty for anonymous users)
	var requestingUserID string
	if user, ok := context.GetUser(r.Context()); ok {
		requestingUserID = user.ID
	}

	// Get followers
	followers, err := h.Service.GetFollowers(userID, requestingUserID)
	if err != nil {
		log.Printf("Error getting followers: %v", err)
		if strings.Contains(err.Error(), "not authorized") {
			http.Error(w, "Not authorized to view followers", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to get followers", http.StatusInternalServerError)
		}
		return
	}

	// Return response
	response := model.FollowersResponse{
		Status: "success",
		Data:   followers,
		Count:  len(followers),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFollowing handles GET /users/:id/following
func (h *FollowersHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path
	userID := extractUserIDFromPath(r.URL.Path, "following")
	if userID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Handle "me" by resolving to current user's ID
	if userID == "me" {
		if user, ok := context.GetUser(r.Context()); ok {
			userID = user.ID
		} else {
			http.Error(w, "Authentication required to view your own following list", http.StatusUnauthorized)
			return
		}
	}

	// Get requesting user ID from context (may be empty for anonymous users)
	var requestingUserID string
	if user, ok := context.GetUser(r.Context()); ok {
		requestingUserID = user.ID
	}

	// Get following
	following, err := h.Service.GetFollowing(userID, requestingUserID)
	if err != nil {
		log.Printf("Error getting following: %v", err)
		if strings.Contains(err.Error(), "not authorized") {
			http.Error(w, "Not authorized to view following list", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to get following list", http.StatusInternalServerError)
		}
		return
	}

	// Return response
	response := model.FollowersResponse{
		Status: "success",
		Data:   following,
		Count:  len(following),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// FollowUser handles POST /users/:id/follow
func (h *FollowersHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUser := context.MustGetUser(r.Context())

	// Extract user ID from URL path
	followedUserID := extractUserIDFromPath(r.URL.Path, "follow")
	if followedUserID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Follow user
	message, err := h.Service.FollowUser(currentUser.ID, followedUserID)
	if err != nil {
		log.Printf("Error following user: %v", err)
		if strings.Contains(err.Error(), "cannot follow yourself") ||
			strings.Contains(err.Error(), "already following") ||
			strings.Contains(err.Error(), "already sent") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": message,
	})
}

// UnfollowUser handles DELETE /users/:id/follow
func (h *FollowersHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUser := context.MustGetUser(r.Context())

	// Extract user ID from URL path
	followedUserID := extractUserIDFromPath(r.URL.Path, "follow")
	if followedUserID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Unfollow user
	err := h.Service.UnfollowUser(currentUser.ID, followedUserID)
	if err != nil {
		log.Printf("Error unfollowing user: %v", err)
		if strings.Contains(err.Error(), "no follow relationship found") {
			http.Error(w, "Not following this user", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Successfully unfollowed user",
	})
}

// GetFollowerCounts handles GET /users/:id/followers/count
func (h *FollowersHandler) GetFollowerCounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path
	userID := extractUserIDFromPath(r.URL.Path, "followers/count")
	if userID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Handle "me" by resolving to current user's ID
	if userID == "me" {
		if user, ok := context.GetUser(r.Context()); ok {
			userID = user.ID
		} else {
			http.Error(w, "Authentication required to view your own follower counts", http.StatusUnauthorized)
			return
		}
	}

	// Get counts
	followersCount, followingCount, err := h.Service.GetFollowerCounts(userID)
	if err != nil {
		log.Printf("Error getting follower counts: %v", err)
		http.Error(w, "Failed to get follower counts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]int{
			"followers": followersCount,
			"following": followingCount,
		},
	})
}

// GetPendingRequests handles GET /users/me/follow-requests
func (h *FollowersHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUser := context.MustGetUser(r.Context())

	// Get pending requests
	requests, err := h.Service.GetPendingFollowRequests(currentUser.ID)
	if err != nil {
		log.Printf("Error getting pending requests: %v", err)
		http.Error(w, "Failed to get pending requests", http.StatusInternalServerError)
		return
	}

	// Return response
	response := model.FollowersResponse{
		Status: "success",
		Data:   requests,
		Count:  len(requests),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AcceptFollowRequest handles POST /users/me/follow-requests/:id/accept
func (h *FollowersHandler) AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUser := context.MustGetUser(r.Context())

	// Extract follower ID from URL path
	followerID := extractFollowerIDFromPath(r.URL.Path)
	if followerID == "" {
		http.Error(w, "Invalid follower ID", http.StatusBadRequest)
		return
	}

	// Accept request
	err := h.Service.AcceptFollowRequest(followerID, currentUser.ID)
	if err != nil {
		log.Printf("Error accepting follow request: %v", err)
		if strings.Contains(err.Error(), "no pending follow request found") {
			http.Error(w, "No pending follow request found", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to accept follow request", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Follow request accepted",
	})
}

// RejectFollowRequest handles POST /users/me/follow-requests/:id/reject
func (h *FollowersHandler) RejectFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUser := context.MustGetUser(r.Context())

	// Extract follower ID from URL path
	followerID := extractFollowerIDFromPath(r.URL.Path)
	if followerID == "" {
		http.Error(w, "Invalid follower ID", http.StatusBadRequest)
		return
	}

	// Reject request
	err := h.Service.RejectFollowRequest(followerID, currentUser.ID)
	if err != nil {
		log.Printf("Error rejecting follow request: %v", err)
		if strings.Contains(err.Error(), "no pending follow request found") {
			http.Error(w, "No pending follow request found", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to reject follow request", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Follow request rejected",
	})
}

// Helper functions

// extractUserIDFromPath extracts user ID from URL paths like /api/users/:id/followers
func extractUserIDFromPath(path, endpoint string) string {
	// Remove leading slash and split by '/'
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	// Expected format: api/users/:id/followers or api/users/:id/following
	if len(parts) >= 4 && parts[0] == "api" && parts[1] == "users" && parts[3] == endpoint {
		return parts[2]
	}

	// For paths like /api/users/:id/follow
	if len(parts) >= 4 && parts[0] == "api" && parts[1] == "users" && parts[3] == endpoint {
		return parts[2]
	}

	// For paths like /api/users/:id/followers/count
	if len(parts) >= 5 && parts[0] == "api" && parts[1] == "users" && parts[3] == "followers" && parts[4] == "count" {
		return parts[2]
	}

	return ""
}

// extractFollowerIDFromPath extracts follower ID from paths like /api/users/me/follow-requests/:id/accept
func extractFollowerIDFromPath(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	// Expected format: api/users/me/follow-requests/:id/accept or api/users/me/follow-requests/:id/reject
	if len(parts) >= 6 && parts[0] == "api" && parts[1] == "users" && parts[2] == "me" && parts[3] == "follow-requests" {
		return parts[4]
	}

	return ""
}

// NewFollowersHandler creates a new followers handler
func NewFollowersHandler(db *sql.DB) *FollowersHandler {
	followerRepo := &repository.FollowerRepository{DB: db}
	userRepo := &repository.UserRepository{DB: db}
	followerService := &service.FollowerService{
		FollowerRepo: followerRepo,
		UserRepo:     userRepo,
	}

	return &FollowersHandler{
		Service: followerService,
	}
}
