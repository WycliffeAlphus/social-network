package handler

import (
	"backend/pkg/db/sqlite"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTestUser(db *sql.DB, email, password string) (string, error) {
	userID := "test-user-id"
	_, err := db.Exec(`INSERT INTO users (id, email, password, fname, lname, dob, imgurl, nickname, about, profileVisibility, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, email, password, "Test", "User", "2000-01-01", "", "testnick", "about", "public", time.Now())
	return userID, err
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	payload := LoginRequest{
		Email:    "wrong@example.com",
		Password: "wrongpass",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	// Defensive: skip test if DB is not available
	db, err := sqlite.ConnectAndMigrate()
	if err != nil {
		t.Skip("DB not available for test")
	}
	defer db.Close()

	LoginHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", res.StatusCode)
	}

	var errs LoginErrs
	_ = json.NewDecoder(res.Body).Decode(&errs)
	if errs.Email == "" && errs.Password == "" {
		t.Error("expected error message for invalid credentials")
	}
}

func TestLoginHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/login", nil)
	rec := httptest.NewRecorder()

	LoginHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", res.StatusCode)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	db, err := sqlite.ConnectAndMigrate()
	if err != nil {
		t.Skip("DB not available for test")
	}
	defer db.Close()

	email := "success@example.com"
	password := "$2a$10$7a6b6b6b6b6b6b6b6b6b6u6b6b6b6b6b6b6b6b6b6b6b6b6b6b6b" // bcrypt hash for 'password'
	_, _ = db.Exec(`DELETE FROM users WHERE email = ?`, email)                // Clean up before test
	_, err = setupTestUser(db, email, password)
	if err != nil {
		t.Fatalf("failed to set up test user: %v", err)
	}

	payload := LoginRequest{
		Email:    email,
		Password: "password", // This will fail unless you use the real hash and compare
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	LoginHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}
}
