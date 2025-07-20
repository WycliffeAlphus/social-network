package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware_NoSessionCookie(t *testing.T) {
	// var db *sql.DB = nil

	// testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("success"))
	// })

	// authHandler := AuthMiddleware(db)(testHandler)
	// req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	// authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}

	expectedBody := "Unauthorized: No session cookie\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, rr.Body.String())
	}
}

func TestAuthMiddleware_EmptySessionCookie(t *testing.T) {
	// var db *sql.DB = nil

	// testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("success"))
	// })

	// authHandler := AuthMiddleware(db)(testHandler)
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "social-network",
		Value: "",
	})

	rr := httptest.NewRecorder()
	// authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}

	expectedBody := "Unauthorized: Empty session token\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, rr.Body.String())
	}
}
