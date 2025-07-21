package handler

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	ctxpkg "backend/internal/context"
	"backend/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

func TestFollowUser_InvalidMethod(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	req := httptest.NewRequest(http.MethodGet, "/api/users/follow", nil)
	// Inject mock user into context
	mockUser := &model.User{ID: "test-user"}
	ctx := ctxpkg.WithUser(req.Context(), mockUser)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	handler := FollowUser(db)
	handler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", recorder.Code)
	}
}

func TestAcceptFollowRequest_InvalidMethod(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	req := httptest.NewRequest(http.MethodGet, "/api/follow/accept", nil)
	mockUser := &model.User{ID: "test-user"}
	ctx := ctxpkg.WithUser(req.Context(), mockUser)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	handler := AcceptFollowRequest(db)
	handler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", recorder.Code)
	}
}

func TestDeclineFollowRequest_InvalidMethod(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	req := httptest.NewRequest(http.MethodGet, "/api/follow/decline", nil)
	mockUser := &model.User{ID: "test-user"}
	ctx := ctxpkg.WithUser(req.Context(), mockUser)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	handler := DeclineFollowRequest(db)
	handler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", recorder.Code)
	}
}

func TestFollowUser_InvalidBody(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	req := httptest.NewRequest(http.MethodPost, "/api/users/follow", bytes.NewBuffer([]byte("notjson")))
	mockUser := &model.User{ID: "test-user"}
	ctx := ctxpkg.WithUser(req.Context(), mockUser)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	handler := FollowUser(db)
	handler.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", recorder.Code)
	}
}
