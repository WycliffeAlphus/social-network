package handler

import (
	"backend/pkg/db/sqlite"
	"net/http"
	"time"
)

// LogoutHandler handles user logout requests by clearing the session cookie and deleting the session from the DB.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST for logout.", http.StatusMethodNotAllowed)
		return
	}

	// Get social-network from cookie
	cookie, err := r.Cookie("social-network")
	if err == nil && cookie.Value != "" {
		// Open DB connection
		db, dbErr := sqlite.ConnectAndMigrate()
		if dbErr == nil {
			defer db.Close()
			_ = sqlite.DeleteSession(db, cookie.Value) // Ignore error for now
		}
	}

	expiresAt := time.Now().In(eat).Add(-1 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "social-network",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
		Expires:  expiresAt,
	})
	http.Error(w, "logout successful", http.StatusOK)
}
