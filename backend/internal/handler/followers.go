package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"backend/pkg/extractid"
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

// AcceptFollowRequest allows the recipient of a follow request to accept it.
// Only the user who is being followed (the recipient) can accept a pending follow request.
// The function checks that the request exists and is pending, then updates its status to 'accepted'.
func AcceptFollowRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get userID from context (this is the user accepting the request)
		currentUserID := context.MustGetUser(r.Context()).ID

		var request struct {
			FollowerID string `json:"followerId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Verify that the follow request exists and is pending
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM followers 
				WHERE follower_id = ? AND followed_id = ? AND status = 'pending'
			)`, request.FollowerID, currentUserID).Scan(&exists)

		if err != nil {
			log.Printf("Error checking follow request: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !exists {
			http.Error(w, "Follow request not found or already processed", http.StatusNotFound)
			return
		}

		// Update the follow request status to accepted
		_, err = db.Exec(`
			UPDATE followers 
			SET status = 'accepted' 
			WHERE follower_id = ? AND followed_id = ? AND status = 'pending'
		`, request.FollowerID, currentUserID)

		if err != nil {
			log.Printf("Error accepting follow request: %v", err)
			http.Error(w, "Failed to accept follow request", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Follow request accepted"})
	}
}

// DeclineFollowRequest allows the recipient of a follow request to decline it.
// Only the user who is being followed (the recipient) can decline a pending follow request.
// The function checks that the request exists and is pending, then deletes it from the database.
func DeclineFollowRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get userID from context (this is the user declining the request)
		currentUserID := context.MustGetUser(r.Context()).ID

		var request struct {
			FollowerID string `json:"followerId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Verify that the follow request exists and is pending
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM followers 
				WHERE follower_id = ? AND followed_id = ? AND status = 'pending'
			)`, request.FollowerID, currentUserID).Scan(&exists)

		if err != nil {
			log.Printf("Error checking follow request: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !exists {
			http.Error(w, "Follow request not found or already processed", http.StatusNotFound)
			return
		}

		// Delete the follow request
		_, err = db.Exec(`
			DELETE FROM followers 
			WHERE follower_id = ? AND followed_id = ? AND status = 'pending'
		`, request.FollowerID, currentUserID)

		if err != nil {
			log.Printf("Error declining follow request: %v", err)
			http.Error(w, "Failed to decline follow request", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Follow request declined"})
	}
}

// GetFollowers handles GET /users/:id/followers
func GetFollowers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserId := context.MustGetUser(r.Context()).ID

		// Extract user ID from URL path
		requestedID := extractid.ExtractUserIDFromPath(r.URL.Path, "followers")
		if requestedID == "" {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// query to get all followers (people who follow the requested user id)
		query := `
			SELECT u.id, u.fname, u.lname, u.imgurl, status
			FROM followers f
			JOIN users u ON f.follower_id = u.id
			WHERE f.followed_id = ? AND f.status = 'accepted'
			ORDER BY f.created_at DESC
		`

		rows, err := db.Query(query, requestedID)
		if err != nil {
			http.Error(w, "Failed to query followers: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var followers []model.UserInfo
		for rows.Next() {
			var user model.UserInfo
			err := rows.Scan(&user.ID, &user.FName, &user.LName, &user.ImgURL, &user.Status)
			if err != nil {
				http.Error(w, "Failed to scan follower: "+err.Error(), http.StatusInternalServerError)
				return
			}
			followers = append(followers, user)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, "Error iterating followers: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// fmt.Println(followers)

		response := model.FollowersResponse{
			Users:         followers,
			CurrentUserId: currentUserId,
			RequestedID:   requestedID,
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetFollowing handles GET /users/:id/following
func GetFollowing(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserId := context.MustGetUser(r.Context()).ID

		// Extract user ID from URL path
		requestedID := extractid.ExtractUserIDFromPath(r.URL.Path, "following")
		if requestedID == "" {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// query to get all users being followed by the requested user
		query := `
    		SELECT u.id, u.fname, u.lname, u.imgurl
    		FROM followers f
    		JOIN users u ON f.followed_id = u.id
    		WHERE f.follower_id = ? AND f.status = 'accepted'
    		ORDER BY f.created_at DESC
		`

		rows, err := db.Query(query, requestedID)
		if err != nil {
			http.Error(w, "Failed to query following: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var following []model.UserInfo
		for rows.Next() {
			var user model.UserInfo
			err := rows.Scan(&user.ID, &user.FName, &user.LName, &user.ImgURL)
			if err != nil {
				http.Error(w, "Failed to scan following user: "+err.Error(), http.StatusInternalServerError)
				return
			}
			following = append(following, user)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, "Error iterating following: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := model.FollowersResponse{
			Users:         following,
			CurrentUserId: currentUserId,
			RequestedID:   requestedID,
		}

		// fmt.Println(response)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetFollowRequests(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserId := context.MustGetUser(r.Context()).ID

		query := `
		SELECT 
			f.follower_id,
			u.fname
			u.lname
			f.imgurl
		FROM followers f
		JOIN users u ON f.follower_id = u.id
		WHERE f.followed_id = ? AND f.status = 'pending'
		ORDER BY f.created_at DESC
	`

		rows, err := db.Query(query, currentUserId)
		if err != nil {
			log.Println("error querying follow requests: ", err)
			http.Error(w, "An error occured. Please check back later", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var requests []model.FollowRequest
		for rows.Next() {
			var req model.FollowRequest
			err := rows.Scan(
				&req.FollowerID,
				&req.FollowerFname,
				&req.FollowerLname,
				&req.FollowerAvatar,
			)
			if err != nil {
				log.Println("error scanning follow request: ", err)
				http.Error(w, "An error occurred processing your request", http.StatusInternalServerError)
				return
			}
			requests = append(requests, req)
		}

		if err = rows.Err(); err != nil {
			log.Println("error after scanning rows: ", err)
			http.Error(w, "An error occurred processing your request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(requests); err != nil {
			log.Println("error encoding response: ", err)
			http.Error(w, "An error occurred processing your request", http.StatusInternalServerError)
		}
	}
}
