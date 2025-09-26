package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"backend/pkg/db/sqlite"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDashboardHandler(t *testing.T) {
	// Defensive: skip test if DB is not available
	db, err := sqlite.ConnectAndMigrate()
	if err != nil {
		t.Skip("DB not available for test")
	}
	defer db.Close()

	// Create test user
	user := &model.User{
		ID:        "test-user-id",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
	}

	// Create request with user in context
	req := httptest.NewRequest("GET", "/api/dashboard", nil)
	ctx := context.WithUser(req.Context(), user)
	req = req.WithContext(ctx)

	// Record response
	rr := httptest.NewRecorder()

	DashboardHandler(db)(rr, req)

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

	// Check welcome message
	expectedMessage := "Welcome to your dashboard, " + user.FirstName + "!"
	if message, ok := data["welcome_message"].(string); !ok || message != expectedMessage {
		t.Errorf("Expected welcome message '%s', got %v", expectedMessage, data["welcome_message"])
	}

	// Check user info
	userInfo, ok := data["user_info"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'user_info' field in response")
	}

	if id, ok := userInfo["id"].(string); !ok || id != user.ID {
		t.Errorf("Expected user ID %s, got %v", user.ID, userInfo["id"])
	}

	if firstName, ok := userInfo["first_name"].(string); !ok || firstName != user.FirstName {
		t.Errorf("Expected first name %s, got %v", user.FirstName, userInfo["first_name"])
	}

	// Check stats
	stats, ok := data["stats"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'stats' field in response")
	}

	// Verify stats are present (even if they're placeholder values)
	if _, exists := stats["posts"]; !exists {
		t.Error("Expected 'posts' in stats")
	}
	if _, exists := stats["followers"]; !exists {
		t.Error("Expected 'followers' in stats")
	}
	if _, exists := stats["following"]; !exists {
		t.Error("Expected 'following' in stats")
	}
}

func TestDashboardHandler_MethodNotAllowed(t *testing.T) {
	// Create POST request (should be rejected)
	req := httptest.NewRequest("POST", "/api/dashboard", nil)

	// Record response
	rr := httptest.NewRecorder()
	DashboardHandler(nil)(rr, req)

	// Check response status
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}
