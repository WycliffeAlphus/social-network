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
