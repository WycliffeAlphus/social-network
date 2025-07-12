package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler_Success(t *testing.T) {
	reqBody := map[string]string{"email": "testuser@example.com", "password": "testpass"}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	LoginHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestLoginHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()
	LoginHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestLoginHandler_BadJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("notjson")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	LoginHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	reqBody := map[string]string{"email": "wrong@example.com", "password": "wrongpass"}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	LoginHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
} 