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
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed. Use GET for dashboard.", http.StatusMethodNotAllowed)
			return
		}

		user := context.MustGetUser(r.Context())

		posts, err := repository.GetPosts(user.ID, db)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		// Get stats
		var postsCount, followersCount, followingCount int
		err = db.QueryRow(`SELECT COUNT(*) FROM posts WHERE user_id = ?`, user.ID).Scan(&postsCount)
		if err != nil {
			postsCount = 0
		}
		err = db.QueryRow(`SELECT COUNT(*) FROM followers WHERE followed_id = ? AND status = 'accepted'`, user.ID).Scan(&followersCount)
		if err != nil {
			followersCount = 0
		}
		err = db.QueryRow(`SELECT COUNT(*) FROM followers WHERE follower_id = ? AND status = 'accepted'`, user.ID).Scan(&followingCount)
		if err != nil {
			followingCount = 0
		}

		// Return dashboard data for the authenticated user
		response := map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"welcome_message": "Welcome to your dashboard, " + user.FirstName + "!",
				"user_info": map[string]interface{}{
					"id":         user.ID,
					"first_name": user.FirstName,
				},
				"stats": map[string]interface{}{
					"posts":     postsCount,
					"followers": followersCount,
					"following": followingCount,
				},
				"posts": posts,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
