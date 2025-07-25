package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"backend/pkg/db/sqlite"
	"backend/pkg/getusers"
	"backend/pkg/models"
	"encoding/json"
	"net/http"
)

// convertModelsUserToInternalUser converts a models.User to internal/model.User
func convertModelsUserToInternalUser(modelsUser models.User) model.User {
	return model.User{
		ID:                modelsUser.ID,
		Email:             modelsUser.Email,
		FirstName:         modelsUser.FirstName,
		LastName:          modelsUser.LastName,
		DOB:               modelsUser.DateOfBirth,
		ImgURL:            modelsUser.AvatarImage.String,
		Nickname:          modelsUser.Nickname.String,
		About:             modelsUser.AboutMe.String,
		ProfileVisibility: modelsUser.ProfileVisibility,
		CreatedAt:         modelsUser.CreatedAt,
	}
}

// ProfileHandler handles requests to get the current user's profile
// This is a protected route that requires authentication
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed. Use GET for profile.", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context (added by auth middleware)
	user := *context.MustGetUser(r.Context())
	userId := user.ID
	db, err := sqlite.ConnectAndMigrate()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	user.ProfileVisibility = "public" // Default to public visibility if its the owner's profile
	id := r.URL.Query().Get("id")
	if id != "" {
		modelsUser, err := getusers.GetUserByID(db, id)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		user = convertModelsUserToInternalUser(modelsUser)
	}
	if user.ProfileVisibility == "private" {
		//check if user is following the profile owner
		following, err := getusers.IsFollowing(db, userId, id)
		if err != nil {
			http.Error(w, "Error checking following status", http.StatusInternalServerError)
			return
		}
		if following {
			user.ProfileVisibility = "public" // Temporarily set to public if following
		}
	}
	// Return user profile (excluding sensitive information like password)
	response := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"id":                 user.ID,
			"email":              user.Email,
			"first_name":         user.FirstName,
			"last_name":          user.LastName,
			"dob":                user.DOB.Format("2006-01-02"),
			"img_url":            user.ImgURL,
			"nickname":           user.Nickname,
			"about":              user.About,
			"profile_visibility": user.ProfileVisibility,
			"created_at":         user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateProfileHandler handles requests to update the current user's profile
// This is a protected route that requires authentication
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed. Use PUT for profile update.", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context (added by auth middleware)
	user := context.MustGetUser(r.Context())
	sessionID := context.MustGetSessionID(r.Context())

	// For now, just return the current user info and session ID
	// In a real implementation, you would parse the request body and update the user
	response := map[string]interface{}{
		"status":  "success",
		"message": "Profile update endpoint (implementation pending)",
		"data": map[string]interface{}{
			"user_id":    user.ID,
			"session_id": sessionID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
