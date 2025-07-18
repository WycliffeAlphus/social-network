package handler

import (
	"backend/pkg/db/sqlite"
	"fmt"
	"log"
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
		log.Printf("Attempting to delete session with ID: '%s'", cookie.Value)
		// Open DB connection
		db, dbErr := sqlite.ConnectAndMigrate()
		if dbErr == nil {
			defer db.Close()
			fmt.Println("cookie value: ", cookie.Value)
			delErr := sqlite.DeleteSession(db, cookie.Value)
			if delErr != nil {
				log.Printf("Failed to delete session from DB: %v", delErr)
			} else {
				log.Printf("Session deleted successfully from DB.")
			}
		}
	}

	expiresAt := time.Now().Add(-1 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "social-network",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to false for local dev
		SameSite: http.SameSiteStrictMode,
		Expires:  expiresAt,
	})
	http.Error(w, "logout successful", http.StatusOK)
}
