package handler

import (
	"backend/internal/auth"
	"encoding/json"
	"net/http"
	"time"
)

var eat = time.FixedZone("EAT", 3*60*60) // East Africa Time (UTC+3)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler handles user login requests, checks credentials, and sets a session cookie.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST for login.", http.StatusMethodNotAllowed)
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format. Please provide email and password.", http.StatusBadRequest)
		return
	}
	// TODO: Replace with real user lookup and bcrypt password check
	// This will be made dynamic when user registration is implemented.
	if req.Email != "testuser@example.com" || req.Password != "testpass" {
		http.Error(w, "Invalid email or password. Please try again.", http.StatusUnauthorized)
		return
	}
	expiresAt := time.Now().In(eat).Add(24 * time.Hour)
	token, err := auth.GenerateJWT(req.Email, 24*time.Hour)
	if err != nil {
		http.Error(w, "Failed to generate authentication token. Please try again.", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  expiresAt,
	})
	http.Error(w, "Login successful. Welcome back!", http.StatusOK)
}
