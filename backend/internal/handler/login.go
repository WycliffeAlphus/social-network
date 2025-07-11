package handler

import (
	"backend/internal/auth"
	"encoding/json"
	"net/http"
	"time"
)

type LoginRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

// LoginHandler handles user login requests, checks credentials, and sets a session cookie.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO: Replace with real user lookup and bcrypt password check
	// This will be made dynamic when user registration is implemented.
	if req.Nickname != "testuser" || req.Password != "testpass" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := auth.GenerateJWT(req.Nickname, time.Hour*24)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	w.Write([]byte("login successful"))
}
