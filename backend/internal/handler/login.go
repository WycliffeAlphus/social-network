package handler

import (
	"backend/internal/auth"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler handles user login requests, checks credentials, and sets a session cookie.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed. Use POST for login.")
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JSON format. Please provide email and password.")
		return
	}
	// TODO: Replace with real user lookup and bcrypt password check
	// This will be made dynamic when user registration is implemented.
	if req.Email != "testuser@example.com" || req.Password != "testpass" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Invalid email or password. Please try again.")
		return
	}
	token, err := auth.GenerateJWT(req.Email, time.Hour*24)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to generate authentication token. Please try again.")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Login successful. Welcome back!")
}
