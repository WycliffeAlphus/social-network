package middlewares

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAuth(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	
	RequireAuth(rr, req)
	
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
	
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}
	
	if status, ok := response["status"].(string); !ok || status != "error" {
		t.Errorf("Expected status 'error', got %v", response["status"])
	}
	
	if message, ok := response["message"].(string); !ok || message != "Authentication required" {
		t.Errorf("Expected message 'Authentication required', got %v", response["message"])
	}
}

func TestAuthMiddleware_NoSessionCookie(t *testing.T) {
	var db *sql.DB = nil
	
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})
	
	authHandler := AuthMiddleware(db)(testHandler)
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	authHandler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
	
	expectedBody := "Unauthorized: No session cookie\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, rr.Body.String())
	}
}

func TestAuthMiddleware_EmptySessionCookie(t *testing.T) {
	var db *sql.DB = nil
	
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})
	
	authHandler := AuthMiddleware(db)(testHandler)
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "social-network",
		Value: "",
	})
	
	rr := httptest.NewRecorder()
	authHandler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
	
	expectedBody := "Unauthorized: Empty session token\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, rr.Body.String())
	}
}

func TestOptionalAuth_NoSessionCookie(t *testing.T) {
	var db *sql.DB = nil
	
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("anonymous"))
	})
	
	optionalAuthHandler := OptionalAuth(db)(testHandler)
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	optionalAuthHandler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	
	if rr.Body.String() != "anonymous" {
		t.Errorf("Expected 'anonymous', got %s", rr.Body.String())
	}
}
