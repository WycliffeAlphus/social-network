package handler

import (
	"backend/internal/context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func GetFollowSuggestions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get userID from context
		currentUserID := context.MustGetSessionID(r.Context())

		// query to get users that the current user is not following
		rows, err := db.Query(`
        SELECT u.id, u.first_name, u.last_name, u.avatar_image,
			EXISTS (
                   SELECT 1 FROM followers f 
                   WHERE f.follower_id = u.id AND f.followed_id = ?
               ) AS follows_me
        FROM users u
        WHERE u.id != ?
        AND NOT EXISTS (
            SELECT 1 FROM followers f 
            WHERE f.follower_id = ? AND f.followed_id = u.id
        )
        LIMIT 7
    `, currentUserID, currentUserID, currentUserID)

		if err != nil {
			log.Printf("Error fetching available users: %v", err)
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []map[string]interface{}
		for rows.Next() {
			var user struct {
				ID        string
				FirstName string
				LastName  string
				Avatar    sql.NullString
				FollowsMe bool
			}

			if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Avatar, &user.FollowsMe); err != nil {
				log.Printf("Error scanning user: %v", err)
				continue
			}

			userMap := map[string]interface{}{
				"id":        user.ID,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"avatar":    user.Avatar,
				"followsMe": user.FollowsMe,
			}

			users = append(users, userMap)
		}

		var visibility string
		err = db.QueryRow("SELECT profile_visibility FROM users WHERE id = ?", currentUserID).Scan(&visibility)
		if err != nil {
			log.Printf("Error fetching privacy setting: %v", err)
			http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"users":      users,
			"visibility": visibility,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
