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

func TestDashboardHandler(t *testing.T) {
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
	DashboardHandler(rr, req)

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
	// Create test user
	user := &model.User{
		ID:    "test-user-id",
		Email: "test@example.com",
	}

	// Create POST request (should be rejected)
	req := httptest.NewRequest("POST", "/api/dashboard", nil)
	ctx := context.WithUser(req.Context(), user)
	req = req.WithContext(ctx)

	// Record response
	rr := httptest.NewRecorder()
	DashboardHandler(rr, req)

	// Check response status
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

func TestPublicDashboardHandler_Authenticated(t *testing.T) {
	// Create test user
	user := &model.User{
		ID:        "test-user-id",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Create request with user in context (authenticated)
	req := httptest.NewRequest("GET", "/api/public-dashboard", nil)
	ctx := context.WithUser(req.Context(), user)
	req = req.WithContext(ctx)

	// Record response
	rr := httptest.NewRecorder()
	PublicDashboardHandler(rr, req)

	// Check response status
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
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

	// Check authenticated flag
	if authenticated, ok := data["authenticated"].(bool); !ok || !authenticated {
		t.Errorf("Expected authenticated to be true, got %v", data["authenticated"])
	}

	// Check personalized message
	expectedMessage := "Welcome back, " + user.FirstName + "!"
	if message, ok := data["message"].(string); !ok || message != expectedMessage {
		t.Errorf("Expected message '%s', got %v", expectedMessage, data["message"])
	}

	// Check user ID is present
	if userID, ok := data["user_id"].(string); !ok || userID != user.ID {
		t.Errorf("Expected user ID %s, got %v", user.ID, data["user_id"])
	}
}

func TestPublicDashboardHandler_NotAuthenticated(t *testing.T) {
	// Create request without user in context (not authenticated)
	req := httptest.NewRequest("GET", "/api/public-dashboard", nil)

	// Record response
	rr := httptest.NewRecorder()
	PublicDashboardHandler(rr, req)

	// Check response status
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
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

	// Check authenticated flag
	if authenticated, ok := data["authenticated"].(bool); !ok || authenticated {
		t.Errorf("Expected authenticated to be false, got %v", data["authenticated"])
	}

	// Check generic message
	expectedMessage := "Welcome to our social network!"
	if message, ok := data["message"].(string); !ok || message != expectedMessage {
		t.Errorf("Expected message '%s', got %v", expectedMessage, data["message"])
	}

	// Check login prompt is present
	if _, exists := data["login_prompt"]; !exists {
		t.Error("Expected 'login_prompt' in response for unauthenticated user")
	}

	// Check user ID is not present
	if _, exists := data["user_id"]; exists {
		t.Error("User ID should not be present for unauthenticated user")
	}
}

func TestPublicDashboardHandler_MethodNotAllowed(t *testing.T) {
	// Create POST request (should be rejected)
	req := httptest.NewRequest("POST", "/api/public-dashboard", nil)

	// Record response
	rr := httptest.NewRecorder()
	PublicDashboardHandler(rr, req)

	// Check response status
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}
