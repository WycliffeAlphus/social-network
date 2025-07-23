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

func TestGetFollowers(t *testing.T) {
	// Create a mock followers handler
	// Note: This is a simplified test - in a real scenario you'd use a test database
	handler := &FollowersHandler{
		Service: nil, // Would be a mock service in real tests
	}

	// Create test user
	user := &model.User{
		ID:        "test-user-id",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Create request with user in context
	req := httptest.NewRequest("GET", "/api/users/user123/followers", nil)
	ctx := context.WithUser(req.Context(), user)
	req = req.WithContext(ctx)

	// Record response
	rr := httptest.NewRecorder()

	// Note: This test would fail without a proper mock service
	// In a real implementation, you'd mock the service layer
	// For now, we're just testing the handler structure

	// Test method validation
	req.Method = "POST"
	handler.GetFollowers(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for POST method, got %d", rr.Code)
	}
}

func TestGetFollowing(t *testing.T) {
	handler := &FollowersHandler{
		Service: nil,
	}

	req := httptest.NewRequest("GET", "/api/users/user123/following", nil)
	rr := httptest.NewRecorder()

	// Test method validation
	req.Method = "POST"
	handler.GetFollowing(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for POST method, got %d", rr.Code)
	}
}

func TestFollowUser_MethodValidation(t *testing.T) {
	handler := &FollowersHandler{
		Service: nil,
	}

	// Test with GET method (should fail)
	req := httptest.NewRequest("GET", "/api/users/user123/follow", nil)
	rr := httptest.NewRecorder()

	handler.FollowUser(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for GET method, got %d", rr.Code)
	}
}

func TestUnfollowUser_MethodValidation(t *testing.T) {
	handler := &FollowersHandler{
		Service: nil,
	}

	// Test with GET method (should fail)
	req := httptest.NewRequest("GET", "/api/users/user123/follow", nil)
	rr := httptest.NewRecorder()

	handler.UnfollowUser(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for GET method, got %d", rr.Code)
	}
}

func TestGetFollowerCounts_MethodValidation(t *testing.T) {
	handler := &FollowersHandler{
		Service: nil,
	}

	// Test with POST method (should fail)
	req := httptest.NewRequest("POST", "/api/users/user123/followers/count", nil)
	rr := httptest.NewRecorder()

	handler.GetFollowerCounts(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for POST method, got %d", rr.Code)
	}
}

func TestExtractUserIDFromPath(t *testing.T) {
	tests := []struct {
		path     string
		endpoint string
		expected string
	}{
		{"/api/users/user123/followers", "followers", "user123"},
		{"/api/users/user456/following", "following", "user456"},
		{"/api/users/user789/follow", "follow", "user789"},
		{"/api/users/user101/followers/count", "followers/count", "user101"},
		{"/invalid/path", "followers", ""},
		{"/api/users/followers", "followers", ""},
	}

	for _, test := range tests {
		result := extractUserIDFromPath(test.path, test.endpoint)
		if result != test.expected {
			t.Errorf("extractUserIDFromPath(%s, %s) = %s, expected %s",
				test.path, test.endpoint, result, test.expected)
		}
	}
}

func TestExtractFollowerIDFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/api/users/me/follow-requests/user123/accept", "user123"},
		{"/api/users/me/follow-requests/user456/reject", "user456"},
		{"/invalid/path", ""},
		{"/api/users/me/follow-requests/accept", ""},
	}

	for _, test := range tests {
		result := extractFollowerIDFromPath(test.path)
		if result != test.expected {
			t.Errorf("extractFollowerIDFromPath(%s) = %s, expected %s",
				test.path, result, test.expected)
		}
	}
}

func TestFollowersResponse_Structure(t *testing.T) {
	// Test the response structure
	followers := []model.FollowerWithUser{
		{
			FollowerID:        "follower1",
			FollowedID:        "user123",
			Status:            "accepted",
			CreatedAt:         time.Now(),
			UserID:            "follower1",
			FirstName:         "John",
			LastName:          "Doe",
			Email:             "john@example.com",
			ProfileVisibility: "public",
		},
	}

	response := model.FollowersResponse{
		Status: "success",
		Data:   followers,
		Count:  len(followers),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled model.FollowersResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.Status != "success" {
		t.Errorf("Expected status 'success', got %s", unmarshaled.Status)
	}

	if unmarshaled.Count != 1 {
		t.Errorf("Expected count 1, got %d", unmarshaled.Count)
	}

	if len(unmarshaled.Data) != 1 {
		t.Errorf("Expected 1 follower, got %d", len(unmarshaled.Data))
	}

	if unmarshaled.Data[0].FirstName != "John" {
		t.Errorf("Expected first name 'John', got %s", unmarshaled.Data[0].FirstName)
	}
}
