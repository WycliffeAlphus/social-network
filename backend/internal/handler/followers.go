package handler

import (
	"backend/internal/context"
	"backend/internal/repository"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// GetFollowers handles GET /users/:id/followers
func GetFollowers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get current user from context (authentication required)
		currentUser := context.MustGetUser(r.Context())

		// Extract user ID from URL path
		userID := extractUserIDFromPath(r.URL.Path, "followers")
		if userID == "" {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Handle "me" by resolving to current user's ID
		if userID == "me" {
			userID = currentUser.ID
		}

		// Get followers from repository
		followerRepo := &repository.FollowerRepository{DB: db}
		followers, err := followerRepo.GetFollowers(userID)
		if err != nil {
			log.Printf("Error getting followers: %v", err)
			http.Error(w, "Failed to get followers", http.StatusInternalServerError)
			return
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   followers,
			"count":  len(followers),
		})
	}
}

// GetFollowing handles GET /users/:id/following
func GetFollowing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get current user from context (authentication required)
		currentUser := context.MustGetUser(r.Context())

		// Extract user ID from URL path
		userID := extractUserIDFromPath(r.URL.Path, "following")
		if userID == "" {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Handle "me" by resolving to current user's ID
		if userID == "me" {
			userID = currentUser.ID
		}

		// Get following from repository
		followerRepo := &repository.FollowerRepository{DB: db}
		following, err := followerRepo.GetFollowing(userID)
		if err != nil {
			log.Printf("Error getting following: %v", err)
			http.Error(w, "Failed to get following", http.StatusInternalServerError)
			return
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   following,
			"count":  len(following),
		})
	}
}

// extractUserIDFromPath extracts user ID from URL paths like /api/users/:id/followers
func extractUserIDFromPath(path, endpoint string) string {
	// Remove leading slash and split by '/'
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	// Expected format: api/users/:id/followers or api/users/:id/following
	if len(parts) >= 4 && parts[0] == "api" && parts[1] == "users" && parts[3] == endpoint {
		return parts[2]
	}

	return ""
}
