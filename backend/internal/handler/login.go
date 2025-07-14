package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"backend/pkg/db/sqlite"
	"backend/pkg/getusers"
	"backend/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

var eat = time.FixedZone("EAT", 3*60*60) // East Africa Time (UTC+3)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginErrs struct {
	Email    string
	Password string
}

// LoginHandler handles user login requests, checks credentials, and sets a session cookie.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST for login.", http.StatusMethodNotAllowed)
		return
	}

	errs := LoginErrs{}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format. Please provide email and password.", http.StatusBadRequest)
		return
	}

	// Open DB connection (in real code, use a shared DB instance)
	db, err := sqlite.ConnectAndMigrate()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	email := strings.ToLower(req.Email)

	// check if email is in db
	var count int
	emailQuerryErr := db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, email).Scan(&count)
	if emailQuerryErr != nil {
		log.Printf("Error querying database while checking if email exists: %v\n", emailQuerryErr)
	}
	// if not, email not found err
	if count < 1 {
		errs.Email = "We did not find an account with that email"
	}

	var user models.User
	if count > 0 {
		var getUserErr error
		user, getUserErr = getusers.GetUserByEmail(db, email)
		if getUserErr != nil {
			log.Println("An error occured while getting user data ByEmailOrNickname", getUserErr)
		}
	}

	compareHashErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	// if password don't match, incorrect password err
	if compareHashErr != nil {
		errs.Password = "Incorrect password. Try again"
	}

	if count < 1 || compareHashErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	} else {
		expiresAt := time.Now().In(eat).Add(24 * time.Hour)
		sessionID, err := sqlite.InsertSession(db, user.ID, expiresAt)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "social-network",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Expires:  expiresAt,
		})

		response := map[string]interface{}{
			"message": "Login successful",
			"user": map[string]interface{}{
				"id":        user.ID,
				"email":     user.Email,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"avatar":    user.AvatarImage,
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
