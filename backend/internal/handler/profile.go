package handler

import (
	"backend/internal/context"
	"backend/pkg/extractid"
	"backend/pkg/getusers"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// ProfileHandler handles requests to get the current user's profile
// This is a protected route that requires authentication
func ProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed. Use GET for profile.", http.StatusMethodNotAllowed)
			return
		}

		requestedID := extractid.ExtractUserIDFromPath(r.URL.Path, "profile")
		currentUserId := context.MustGetUser(r.Context()).ID

		// handle "current" user special case
		switch requestedID {
		case "current":
			requestedID = currentUserId
		case "":
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := getusers.GetUserByID(db, requestedID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Profile not found", http.StatusNotFound)
				return
			} else {
				log.Println("Error getitng user by id", err)
				http.Error(w, "An error occured, check back later", http.StatusInternalServerError)
				return
			}
		}

		// get follow status
		var followsMe bool
		followsMeQuery := `
				SELECT EXISTS (
                   SELECT 1 FROM followers f 
                   WHERE f.follower_id = ? AND f.followed_id = ?
                ) AS follows_me
		`
		followsMeCheckErr := db.QueryRow(followsMeQuery, currentUserId, requestedID).Scan(&followsMe)
		if followsMeCheckErr != nil {
			if followsMeCheckErr == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			log.Println("Error getting user by id: ", followsMeCheckErr)
			http.Error(w, "An error occured, please check back later", http.StatusInternalServerError)
			return
		}

		// get followers count
		var followersCount int
		err = db.QueryRow(`
		SELECT COUNT(*) FROM followers
		WHERE followed_id = ? AND status = 'accepted'
		`, requestedID).Scan(&followersCount)
		if err != nil {
			log.Println("Error getting followers count:", err)
			followersCount = 0 // default to 0 if there's an error
		}

		// get following count
		var followingCount int
		err = db.QueryRow(`SELECT COUNT(*) FROM followers
		WHERE follower_id = ? AND status = 'accepted'
		`, requestedID).Scan(&followingCount)
		if err != nil {
			log.Println("Error getting following count:", err)
			followingCount = 0 // default to 0 if there's an error
		}

		response := map[string]interface{}{
			"current_user_id": currentUserId,
			"follows_me":      followsMe,
			"followers_count": followersCount,
			"following_count": followingCount,
			"profile": map[string]interface{}{
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
}
