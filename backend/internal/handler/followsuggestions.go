package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func GetFollowSuggestions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get userID from context
		currentUserID := r.Context().Value("userID").(string)

		// query to get users that the current user is not following
		rows, err := db.Query(`
        SELECT u.id, u.first_name, u.last_name, u.avatar_image
        FROM users u
        WHERE u.id != ?
        AND NOT EXISTS (
            SELECT 1 FROM followers f 
            WHERE f.follower_id = ? AND f.followed_id = u.id
        )
        LIMIT 7
    `, currentUserID, currentUserID)

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
			}

			if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Avatar); err != nil {
				log.Printf("Error scanning user: %v", err)
				continue
			}

			userMap := map[string]interface{}{
				"id":        user.ID,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"avatar":    user.Avatar,
			}

			users = append(users, userMap)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
