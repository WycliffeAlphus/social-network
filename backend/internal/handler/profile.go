package handler

import (
	"backend/internal/context"
	"encoding/json"
	"net/http"
)

// ProfileHandler handles requests to get the current user's profile
// This is a protected route that requires authentication
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed. Use GET for profile.", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context (added by auth middleware)
	user := context.MustGetUser(r.Context())

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
