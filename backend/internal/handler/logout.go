package handler

import (
	"net/http"
	"time"
)

// LogoutHandler handles user logout requests by clearing the session cookie.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST for logout.", http.StatusMethodNotAllowed)
		return
	}
	expiresAt := time.Now().In(eat).Add(-1 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
		Expires:  expiresAt,
	})
	http.Error(w, "logout successful", http.StatusOK)
}
