package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"backend/internal/service" // Import service package
	"backend/internal/utils"
	"backend/pkg/extractid"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// FollowerHandler holds the services needed by the follower handlers.

type FollowerHandler struct {
	DB                  *sql.DB // Keeping DB for now, can be refactored to a FollowerService later
	NotificationService *service.NotificationService
}

// NewFollowerHandler creates a new FollowerHandler.
func NewFollowerHandler(db *sql.DB, ns *service.NotificationService) *FollowerHandler {
	return &FollowerHandler{
		DB:                  db,
		NotificationService: ns,
	}
}

func (h *FollowerHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// get userID from context
	currentUserID := context.MustGetUser(r.Context()).ID

	var request struct {
		FollowedUserID string `json:"userId"`
		IsFollowBack   bool   `json:"isFollowBack"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// prevent self-follow
	if currentUserID == request.FollowedUserID {
		utils.RespondWithError(w, http.StatusBadRequest, "Cannot follow yourself")
		return
	}

	// fetch profileVisibility of the followed user
	var profileVisibility string
	err := h.DB.QueryRow("SELECT profileVisibility FROM users WHERE id = ?", request.FollowedUserID).Scan(&profileVisibility)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// determine status based on visibility
	status := "requested"
	if profileVisibility == "public" {
		status = "accepted"
	}
	if request.IsFollowBack {
		status = "accepted"
	}

	// check if already following
	var alreadyFollowing bool
	checkAlreadyFollowingErr := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = ? AND followed_id = ?)",
		currentUserID, request.FollowedUserID).Scan(&alreadyFollowing)
	if checkAlreadyFollowingErr != nil {
		log.Printf("Error checking follow status: %v", checkAlreadyFollowingErr)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if alreadyFollowing {
		utils.RespondWithError(w, http.StatusBadRequest, "Already following this user")
		return
	}

	// create follow relationship
	_, insertFollowErr := h.DB.Exec("INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)", currentUserID, request.FollowedUserID, status)
	if insertFollowErr != nil {
		log.Printf("Error creating follow relationship: %v", insertFollowErr)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to follow user")
		return
	}

	// Create a notification
	actorID := currentUserID
	followedUserID := request.FollowedUserID
	if status == "requested" {
		err := h.NotificationService.CreateFollowRequestNotification(actorID, followedUserID)
		if err != nil {
			log.Printf("Error creating follow request notification: %v", err)
			// Non-critical error, so we don't block the user response
		}
	} else if status == "accepted" {
		if request.IsFollowBack {
			err := h.NotificationService.CreateFollowBackNotification(actorID, followedUserID)
			if err != nil {
				log.Printf("Error creating follow back notification: %v", err)
			}
		} else {
			err := h.NotificationService.CreateNewFollowerNotification(actorID, followedUserID)
			if err != nil {
				log.Printf("Error creating new follower notification: %v", err)
				// Non-critical error, so we don't block the user response
			}
		}
	}

	message := "Follow request sent"
	if status == "accepted" {
		message = "Successfully followed user"
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": message,
		"status":  status,
	})
}

// AcceptFollowRequest allows the recipient of a follow request to accept it.
// Only the user who is being followed (the recipient) can accept a pending follow request.
// The function checks that the request exists and is pending, then updates its status to 'accepted'.
func (h *FollowerHandler) AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// get userID from context (this is the user accepting the request)
	currentUserID := context.MustGetUser(r.Context()).ID

	var request struct {
		FollowerID string `json:"followerId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Verify that the follow request exists and is pending
	var exists bool
	err := h.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM followers 
			WHERE follower_id = ? AND followed_id = ? AND status = 'requested'
		)`, request.FollowerID, currentUserID).Scan(&exists)

	if err != nil {
		log.Printf("Error checking follow request: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Follow request not found or already processed")
		return
	}

	// Update the follow request status to accepted
	result, err := h.DB.Exec(`
		UPDATE followers 
		SET status = 'accepted' 
		WHERE follower_id = ? AND followed_id = ? AND status = 'requested'
	`, request.FollowerID, currentUserID)

	if err != nil {
		log.Printf("Error accepting follow request: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to accept follow request")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to accept follow request")
		return
	}

	if rowsAffected == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Follow request not found or already processed")
		return
	}

	if err != nil {
		log.Printf("Error accepting follow request: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to accept follow request")
		return
	}

	// Create a notification for the user who sent the follow request
	followerID := request.FollowerID
	actorID := currentUserID
	if err := h.NotificationService.CreateFollowAcceptedNotification(actorID, followerID); err != nil {
		log.Printf("Error creating follow accepted notification: %v", err)
		// Non-critical error, so we don't block the user response
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Follow request accepted", "status": "accepted"})
}

// DeclineFollowRequest allows the recipient of a follow request to decline it.
// Only the user who is being followed (the recipient) can decline a pending follow request.
// The function checks that the request exists and is pending, then deletes it from the database.
func (h *FollowerHandler) DeclineFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// get userID from context (this is the user declining the request)
	currentUserID := context.MustGetUser(r.Context()).ID

	var request struct {
		FollowerID string `json:"followerId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Verify that the follow request exists and is pending
	var exists bool
	err := h.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM followers 
			WHERE follower_id = ? AND followed_id = ? AND status = 'requested'
		)`, request.FollowerID, currentUserID).Scan(&exists)

	if err != nil {
		log.Printf("Error checking follow request: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Follow request not found or already processed")
		return
	}

	// Delete the follow request
	_, err = h.DB.Exec(`
		DELETE FROM followers 
		WHERE follower_id = ? AND followed_id = ? AND status = 'requested'
	`, request.FollowerID, currentUserID)

	if err != nil {
		log.Printf("Error declining follow request: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to decline follow request")
		return
	}

	// Create a notification for the user who sent the follow request
	followerID := request.FollowerID
	actorID := currentUserID
	if err := h.NotificationService.CreateFollowDeclinedNotification(actorID, followerID); err != nil {
		log.Printf("Error creating follow declined notification: %v", err)
		// Non-critical error, so we don't block the user response
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Follow request declined"})
}

// CancelFollowRequest allows a user to cancel their own follow request.
// Only the user who sent the request can cancel it.
func (h *FollowerHandler) CancelFollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// get userID from context (this is the user canceling the request)
	currentUserID := context.MustGetUser(r.Context()).ID

	var request struct {
		UserID string `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Verify that the follow request exists and is requested
	var exists bool
	err := h.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM followers 
			WHERE follower_id = ? AND followed_id = ? AND status = 'requested'
		)`, currentUserID, request.UserID).Scan(&exists)

	if err != nil {
		log.Printf("Error checking follow request: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Follow request not found or already processed")
		return
	}

	// Delete the follow request
	_, err = h.DB.Exec(`
		DELETE FROM followers 
		WHERE follower_id = ? AND followed_id = ? AND status = 'requested'
	`, currentUserID, request.UserID)

	if err != nil {
		log.Printf("Error canceling follow request: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to cancel follow request")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Follow request canceled"})
}

// GetFollowers handles GET /users/:id/followers
func (h *FollowerHandler) GetFollowers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		requestedID := extractid.ExtractUserIDFromPath(r.URL.Path, "followers")
		currentUserId := context.MustGetUser(r.Context()).ID

		// handle "current" user special case
		switch requestedID {
		case "currentuser":
			requestedID = currentUserId
		case "":
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
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

		rows, err := h.DB.Query(query, requestedID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to query followers: "+err.Error())
			return
		}
		defer rows.Close()

		var followers []model.UserInfo
		for rows.Next() {
			var user model.UserInfo
			err := rows.Scan(&user.ID, &user.FName, &user.LName, &user.ImgURL, &user.Status)
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to scan follower: "+err.Error())
				return
			}
			followers = append(followers, user)
		}
		if err = rows.Err(); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating followers: "+err.Error())
			return
		}

		// fmt.Println(followers)

		response := model.FollowersResponse{
			Users:         followers,
			CurrentUserId: currentUserId,
			RequestedID:   requestedID,
		}

		// Return response
		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

// GetFollowing handles GET /users/:id/following
func (h *FollowerHandler) GetFollowing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		currentUserId := context.MustGetUser(r.Context()).ID

		// Extract user ID from URL path
		requestedID := extractid.ExtractUserIDFromPath(r.URL.Path, "following")
		if requestedID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
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

		rows, err := h.DB.Query(query, requestedID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to query following: "+err.Error())
			return
		}
		defer rows.Close()

		var following []model.UserInfo
		for rows.Next() {
			var user model.UserInfo
			err := rows.Scan(&user.ID, &user.FName, &user.LName, &user.ImgURL)
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to scan following user: "+err.Error())
				return
			}
			following = append(following, user)
		}
		if err = rows.Err(); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating following: "+err.Error())
			return
		}

		response := model.FollowersResponse{
			Users:         following,
			CurrentUserId: currentUserId,
			RequestedID:   requestedID,
		}

		// fmt.Println(response)

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

// GetFollowRequests handles GET /users/:id/follow-requests
func (h *FollowerHandler) GetFollowRequests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		currentUserId := context.MustGetUser(r.Context()).ID

		query := `
		SELECT 
			f.follower_id,
			u.fname,
			u.lname,
			u.imgurl
		FROM followers f
		JOIN users u ON f.follower_id = u.id
		WHERE f.followed_id = ? AND f.status = 'requested'
		ORDER BY f.created_at DESC
	`

		rows, err := h.DB.Query(query, currentUserId)
		if err != nil {
			log.Println("error querying follow requests: ", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "An error occured. Please check back later")
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
				utils.RespondWithError(w, http.StatusInternalServerError, "An error occurred processing your request")
				return
			}
			requests = append(requests, req)
		}

		if err = rows.Err(); err != nil {
			log.Println("error after scanning rows: ", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "An error occurred processing your request")
			return
		}

		if requests == nil {
			requests = []model.FollowRequest{}
		}

		utils.RespondWithJSON(w, http.StatusOK, requests)
	}
}

// GetFollowStatus checks the follow status between two users
func (h *FollowerHandler) GetFollowStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		currentUserID := context.MustGetUser(r.Context()).ID

		// Extract user ID from URL path
		requestedID := extractid.ExtractUserIDFromPath(r.URL.Path, "follow-status")
		if requestedID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		// Check if current user is following the requested user
		var status string
		err := h.DB.QueryRow(`
			SELECT status 
			FROM followers 
			WHERE follower_id = ? AND followed_id = ?
		`, currentUserID, requestedID).Scan(&status)

		if err != nil {
			if err == sql.ErrNoRows {
				status = "not_following"
			} else {
				log.Printf("Error checking follow status: %v", err)
				utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
		}

		response := map[string]string{
			"status": status,
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

// GetFollowStatuses checks the follow status for a list of users.
func (h *FollowerHandler) GetFollowStatuses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	currentUserID := context.MustGetUser(r.Context()).ID

	var request struct {
		UserIDs []string `json:"userIds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	statuses := make(map[string]string)
	for _, userID := range request.UserIDs {
		var status string
		err := h.DB.QueryRow(`
			SELECT status 
			FROM followers 
			WHERE follower_id = ? AND followed_id = ?
		`, currentUserID, userID).Scan(&status)

		if err != nil {
			if err == sql.ErrNoRows {
				statuses[userID] = "not_following"
			} else {
				log.Printf("Error checking follow status: %v", err)
				statuses[userID] = "error"
			}
		} else {
			statuses[userID] = status
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, statuses)
}

// GetIncomingFollowRequestStatus checks if the viewed user has sent a follow request to the current user.
func (h *FollowerHandler) GetIncomingFollowRequestStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		currentUserID := context.MustGetUser(r.Context()).ID

		// Extract user ID from URL path (this is the ID of the user whose profile is being viewed)
		viewedUserID := extractid.ExtractUserIDFromPath(r.URL.Path, "incoming-follow-request-status")
		if viewedUserID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		// Check if the viewed user has sent a follow request to the current user
		var status string
		err := h.DB.QueryRow(`
			SELECT status 
			FROM followers 
			WHERE follower_id = ? AND followed_id = ?
		`, viewedUserID, currentUserID).Scan(&status)

		if err != nil {
			if err == sql.ErrNoRows {
				status = "none" // No incoming request or already declined/accepted and deleted
			} else {
				log.Printf("Error checking incoming follow request status: %v", err)
				utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
				return
			}
		}

		response := map[string]string{
			"status": status,
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}
