package handler

import (
	"net/http"
	"time"
)

// LogoutHandler handles user logout requests by clearing the session cookie.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("logout successful"))
} 