package middlewares

import (
	"backend/internal/context"
	"backend/internal/model"
	"backend/pkg/db/sqlite"
	"backend/pkg/getusers"
	"database/sql"
	"log"
	"net/http"
	"time"
)

// AuthMiddleware verifies session token and attaches user to context
func AuthMiddleware(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		cookie, err := r.Cookie("social-network")
		if err != nil {
			http.Error(w, "Unauthorized: No session cookie", http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value
		if sessionID == "" {
			http.Error(w, "Unauthorized: Empty session token", http.StatusUnauthorized)
			return
		}

		// Get session from database
		session, err := sqlite.GetSession(db, sessionID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			} else {
				log.Printf("Error retrieving session: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Check if session is expired
		if time.Now().After(session.ExpiresAt) {
			// Clean up expired session
			_ = sqlite.DeleteSession(db, sessionID)
			http.Error(w, "Unauthorized: Session expired", http.StatusUnauthorized)
			return
		}

		// Get user from database
		user, err := getusers.GetUserByID(db, session.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				// User was deleted but session still exists
				_ = sqlite.DeleteSession(db, sessionID)
				http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
			} else {
				log.Printf("Error retrieving user: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Convert pkg/models.User to internal/model.User
		modelUser := &model.User{
			ID:                user.ID,
			Email:             user.Email,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			DOB:               user.DateOfBirth,
			ImgURL:            user.AvatarImage.String,
			Nickname:          user.Nickname.String,
			About:             user.AboutMe.String,
			ProfileVisibility: user.ProfileVisibility,
			CreatedAt:         user.CreatedAt,
		}

		// Add user and session ID to context
		ctx := context.WithUser(r.Context(), modelUser)
		ctx = context.WithSessionID(ctx, session.UserID)

		// Continue with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
