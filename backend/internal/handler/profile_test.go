package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfileHandler_MethodNotAllowed(t *testing.T) {
	// Create POST request (should be rejected)
	req := httptest.NewRequest("POST", "/api/profile", nil)

	// Record response
	rr := httptest.NewRecorder()
	handler := ProfileHandler(nil)
	handler(rr, req)

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

	handler := ProfileHandler(nil)
	handler(rr, req)
}
