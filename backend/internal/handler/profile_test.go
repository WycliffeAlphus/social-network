package handler

import (
	"backend/internal/context"
	"backend/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	// ProfileHandler(rr, req)

	// Check response status
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

func TestProfileHandler_NoUserInContext(t *testing.T) {
	// Create request without user in context
	// req := httptest.NewRequest("GET", "/api/profile", nil)

	// Record response
	// rr := httptest.NewRecorder()

	// This should panic because MustGetUser is called without user in context
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected ProfileHandler to panic when no user in context")
		}
	}()

	// ProfileHandler(rr, req)
}
