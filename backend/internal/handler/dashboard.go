package handler

import (
	"backend/internal/context"
	"backend/internal/repository"
	"database/sql"
	"encoding/json"
	"net/http"
)

func DashboardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := context.MustGetUser(r.Context())

		posts, err := repository.GetPosts(user.ID, db)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		// Return dashboard data for the authenticated user
		response := map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"posts": posts,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
