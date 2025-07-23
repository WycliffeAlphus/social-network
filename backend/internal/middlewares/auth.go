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
			// Convert pkg/models.User to internal/model.User
			modelUser := &model.User{
				ID:                user.ID,
				Email:             user.Email,
				FirstName:         user.FirstName,
				LastName:          user.LastName,
				DOB:               user.DOB,
				ImgURL:            user.ImgURL,
				Nickname:          user.Nickname,
				About:             user.About,
				ProfileVisibility: user.ProfileVisibility,
				CreatedAt:         user.CreatedAt,
			}

		// Add user and session ID to context
		ctx := context.WithUser(r.Context(), modelUser)
		ctx = context.WithSessionID(ctx, sessionID)

		// Continue with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
			// Continue with the request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth is a convenience function that returns a 401 JSON response
func RequireAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": "Authentication required",
	})
}

// OptionalAuth middleware that adds user to context if authenticated, but doesn't require it
func OptionalAuth(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get session cookie
			cookie, err := r.Cookie("social-network")
			if err != nil || cookie.Value == "" {
				// No session, continue without user context
				next.ServeHTTP(w, r)
				return
			}

			sessionID := cookie.Value

			// Get session from database
			session, err := sqlite.GetSession(db, sessionID)
			if err != nil {
				// Invalid session, continue without user context
				next.ServeHTTP(w, r)
				return
			}

			// Check if session is expired
			if time.Now().After(session.ExpiresAt) {
				// Clean up expired session and continue without user context
				_ = sqlite.DeleteSession(db, sessionID)
				next.ServeHTTP(w, r)
				return
			}

			// Get user from database
			user, err := getusers.GetUserByID(db, session.UserID)
			if err != nil {
				// User not found, continue without user context
				next.ServeHTTP(w, r)
				return
			}

			// Convert pkg/models.User to internal/model.User
			modelUser := &model.User{
				ID:                user.ID,
				Email:             user.Email,
				FirstName:         user.FirstName,
				LastName:          user.LastName,
				DOB:               user.DOB,
				ImgURL:            user.ImgURL,
				Nickname:          user.Nickname,
				About:             user.About,
				ProfileVisibility: user.ProfileVisibility,
				CreatedAt:         user.CreatedAt,
			}

			// Add user and session ID to context
			ctx := context.WithUser(r.Context(), modelUser)
			ctx = context.WithSessionID(ctx, sessionID)

			// Continue with the request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
