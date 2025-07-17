package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProfileHandler(t *testing.T) {
	// Create test user
	user := &model.User{
		ID:                "test-user-id",
		Email:             "test@example.com",
		FirstName:         "John",
		LastName:          "Doe",
		DOB:               time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		ImgURL:            "https://example.com/avatar.jpg",
		Nickname:          "johndoe",
		About:             "Test user bio",
		ProfileVisibility: "public",
		CreatedAt:         time.Now(),
	}

	// Create request with user in context
	req := httptest.NewRequest("GET", "/api/profile", nil)
	ctx := context.WithUser(req.Context(), user)
	req = req.WithContext(ctx)

	// Record response
	rr := httptest.NewRecorder()
	ProfileHandler(rr, req)

	// Check response status
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Check content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}

	// Parse response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Check response structure
	if status, ok := response["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %v", response["status"])
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'data' field in response")
	}

	// Check user data in response
	if id, ok := data["id"].(string); !ok || id != user.ID {
		t.Errorf("Expected user ID %s, got %v", user.ID, data["id"])
	}

	if email, ok := data["email"].(string); !ok || email != user.Email {
		t.Errorf("Expected email %s, got %v", user.Email, data["email"])
	}

	if firstName, ok := data["first_name"].(string); !ok || firstName != user.FirstName {
		t.Errorf("Expected first name %s, got %v", user.FirstName, data["first_name"])
	}

	// Verify password is not included in response
	if _, exists := data["password"]; exists {
		t.Error("Password should not be included in profile response")
	}
}

func TestProfileHandler_MethodNotAllowed(t *testing.T) {
	// Create test user
	user := &model.User{
		ID:    "test-user-id",
		Email: "test@example.com",
	}

	// Create POST request (should be rejected)
	req := httptest.NewRequest("POST", "/api/profile", nil)
	ctx := context.WithUser(req.Context(), user)
	req = req.WithContext(ctx)

	// Record response
	rr := httptest.NewRecorder()
	ProfileHandler(rr, req)

	// Check response status
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

func TestProfileHandler_NoUserInContext(t *testing.T) {
	// Create request without user in context
	req := httptest.NewRequest("GET", "/api/profile", nil)

	// Record response
	rr := httptest.NewRecorder()

	// This should panic because MustGetUser is called without user in context
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected ProfileHandler to panic when no user in context")
		}
	}()

	ProfileHandler(rr, req)
}

func TestUpdateProfileHandler(t *testing.T) {
	// Create test user
	user := &model.User{
		ID:    "test-user-id",
		Email: "test@example.com",
	}

	sessionID := "test-session-id"

	// Create request with user and session in context
	req := httptest.NewRequest("PUT", "/api/profile/update", nil)
	ctx := context.WithUser(req.Context(), user)
	ctx = context.WithSessionID(ctx, sessionID)
	req = req.WithContext(ctx)

	// Record response
	rr := httptest.NewRecorder()
	UpdateProfileHandler(rr, req)

	// Check response status
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Check content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}

	// Parse response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Check response structure
	if status, ok := response["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %v", response["status"])
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'data' field in response")
	}

	// Check that user ID and session ID are in response
	if userID, ok := data["user_id"].(string); !ok || userID != user.ID {
		t.Errorf("Expected user ID %s, got %v", user.ID, data["user_id"])
	}

	if respSessionID, ok := data["session_id"].(string); !ok || respSessionID != sessionID {
		t.Errorf("Expected session ID %s, got %v", sessionID, data["session_id"])
	}
}

func TestUpdateProfileHandler_MethodNotAllowed(t *testing.T) {
	// Create test user
	user := &model.User{
		ID:    "test-user-id",
		Email: "test@example.com",
	}

	// Create GET request (should be rejected)
	req := httptest.NewRequest("GET", "/api/profile/update", nil)
	ctx := context.WithUser(req.Context(), user)
	req = req.WithContext(ctx)

	// Record response
	rr := httptest.NewRecorder()
	UpdateProfileHandler(rr, req)

	// Check response status
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

func TestUpdateProfileHandler_NoUserInContext(t *testing.T) {
	// Create request without user in context
	req := httptest.NewRequest("PUT", "/api/profile/update", nil)

	// Record response
	rr := httptest.NewRecorder()

	// This should panic because MustGetUser is called without user in context
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected UpdateProfileHandler to panic when no user in context")
		}
	}()

	UpdateProfileHandler(rr, req)
}
