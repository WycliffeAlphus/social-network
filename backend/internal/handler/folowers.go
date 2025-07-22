package handler

import (
	"backend/internal/context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func FollowUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get userID from context
		currentUserID := context.MustGetUser(r.Context()).ID

		var request struct {
			FollowedUserID string `json:"userId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// prevent self-follow
		if currentUserID == request.FollowedUserID {
			http.Error(w, "Cannot follow yourself", http.StatusBadRequest)
			return
		}

		// fetch profileVisibility of the followed user
		var profileVisibility string
		err := db.QueryRow("SELECT profileVisibility FROM users WHERE id = ?", request.FollowedUserID).Scan(&profileVisibility)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// determine status based on visibility
		status := "pending"
		if profileVisibility == "public" {
			status = "accepted"
		}

		// check if already following
		var alreadyFollowing bool
		checkAlreadyFollowingErr := db.QueryRow("SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = ? AND followed_id = ?)",
			currentUserID, request.FollowedUserID).Scan(&alreadyFollowing)
		if checkAlreadyFollowingErr != nil {
			log.Printf("Error checking follow status: %v", checkAlreadyFollowingErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if alreadyFollowing {
			http.Error(w, "Already following this user", http.StatusBadRequest)
			return
		}

		// create follow relationship
		_, insertFollowErr := db.Exec("INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)", currentUserID, request.FollowedUserID, status)
		if insertFollowErr != nil {
			log.Printf("Error creating follow relationship: %v", insertFollowErr)
			http.Error(w, "Failed to follow user", http.StatusInternalServerError)
			return
		}

		message := "Follow request sent"
		if status == "accepted" {
			message = "Successfully followed user"
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": message})
	}
}
