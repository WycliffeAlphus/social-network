package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoutHandler_Success(t *testing.T) {
	req := httptest.NewRequest("POST", "/logout", nil)
	w := httptest.NewRecorder()
	LogoutHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestLogoutHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("GET", "/logout", nil)
	w := httptest.NewRecorder()
	LogoutHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", resp.StatusCode)
	}
}
