package handler

import (
	"backend/internal/context"
	"encoding/json"
	"net/http"
)

// DashboardHandler handles requests to get dashboard data for the authenticated user
// This is a protected route that requires authentication
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed. Use GET for dashboard.", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context (added by auth middleware)
	user := context.MustGetUser(r.Context())

	// Return dashboard data for the authenticated user
	response := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"welcome_message": "Welcome to your dashboard, " + user.FirstName + "!",
			"user_info": map[string]interface{}{
				"id":         user.ID,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
			},
			"stats": map[string]interface{}{
				"posts":     0, // Placeholder - would be fetched from database
				"followers": 0, // Placeholder - would be fetched from database
				"following": 0, // Placeholder - would be fetched from database
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PublicDashboardHandler demonstrates optional authentication
// This route can be accessed by both authenticated and non-authenticated users
func PublicDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed. Use GET for public dashboard.", http.StatusMethodNotAllowed)
		return
	}

	// Try to get user from context (may or may not be present)
	user, isAuthenticated := context.GetUser(r.Context())

	var response map[string]interface{}

	if isAuthenticated {
		// User is logged in, show personalized content
		response = map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"message":        "Welcome back, " + user.FirstName + "!",
				"authenticated":  true,
				"user_id":        user.ID,
				"public_content": "Here's some public content for logged-in users",
			},
		}
	} else {
		// User is not logged in, show public content only
		response = map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"message":        "Welcome to our social network!",
				"authenticated":  false,
				"public_content": "Here's some public content for everyone",
				"login_prompt":   "Please log in to see personalized content",
			},
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
