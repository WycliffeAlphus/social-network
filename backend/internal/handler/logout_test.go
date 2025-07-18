package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoutHandler_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	rec := httptest.NewRecorder()

	// Set a social-network cookie
	req.AddCookie(&http.Cookie{
		Name:  "social-network",
		Value: "sometoken",
	})

	LogoutHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	// Check that the social-network cookie is cleared
	found := false
	for _, c := range res.Cookies() {
		if c.Name == "social-network" && c.Value == "" {
			found = true
		}
	}
	if !found {
		t.Error("expected social-network cookie to be cleared")
	}
}

func TestLogoutHandler_NoCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	rec := httptest.NewRecorder()

	LogoutHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}
}

func TestLogoutHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/logout", nil)
	rec := httptest.NewRecorder()

	LogoutHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", res.StatusCode)
	}
}
