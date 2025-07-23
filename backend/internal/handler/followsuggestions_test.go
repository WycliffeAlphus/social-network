// followsuggestions_test.go
//
// This test file verifies the GetFollowSuggestions handler.
// It ensures the handler returns HTTP 200 when a user is present in the context and the test database is properly set up.
// The test sets up in-memory users and followers tables and injects a mock user into the request context.

package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	ctxpkg "backend/internal/context"
	"backend/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetFollowSuggestions_WithUserInContext(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	// Create mock tables
	db.Exec(`CREATE TABLE users (id TEXT PRIMARY KEY, fname TEXT, lname TEXT, imgurl TEXT, profileVisibility TEXT)`)
	db.Exec(`CREATE TABLE followers (follower_id TEXT, followed_id TEXT, status TEXT)`)
	db.Exec(`INSERT INTO users (id, fname, lname, imgurl, profileVisibility) VALUES ('test-user', 'Test', 'User', '', 'public')`)

	// Create a request and inject a user into the context
	req := httptest.NewRequest(http.MethodGet, "/api/users/available", nil)
	mockUser := &model.User{ID: "test-user"}
	ctx := ctxpkg.WithUser(req.Context(), mockUser)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()
	handler := GetFollowSuggestions(db)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", recorder.Code)
	}
}
