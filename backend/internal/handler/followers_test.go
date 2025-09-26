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

// MockNotificationService is a mock implementation of NotificationService for testing.
type MockNotificationService struct{}

func (m *MockNotificationService) CreateFollowRequestNotification(actorID, targetUserID string) error {
	return nil
}

func (m *MockNotificationService) CreateFollowAcceptedNotification(actorID, targetUserID string) error {
	return nil
}

func (m *MockNotificationService) CreateNewFollowerNotification(actorID, targetUserID string) error {
	return nil
}

func (m *MockNotificationService) CreateGroupInviteNotification(actorID, targetUserID string, groupID, invitationID int) error {
	return nil
}

func (m *MockNotificationService) GetByUserID(userID string) ([]*model.Notification, error) {
	return nil, nil
}

func (m *MockNotificationService) MarkAllAsRead(userID string) error {
	return nil
}

func (m *MockNotificationService) CreateGroupJoinRequestNotification(actorID, groupOwnerID string, groupID int) error {
	return nil
}

func (m *MockNotificationService) CreateGroupJoinAcceptedNotification(actorID, targetUserID string, groupID int) error {
	return nil
}

func (m *MockNotificationService) CreateGroupEventNotification(actorID string, groupID, eventID int) error {
	return nil
}

func (m *MockNotificationService) CreatePostNotification(actorID, postID string, groupID *int) error {
	return nil
}

func (m *MockNotificationService) CreateCommentNotification(actorID, postOwnerID, postID string) error {
	return nil
}

func (m *MockNotificationService) CreateReactionNotification(actorID, postOwnerID, postID string) error {
	return nil
}

func (m *MockNotificationService) CreateFollowBackNotification(actorID, targetUserID string) error {
	return nil
}

func (m *MockNotificationService) CreateFollowDeclinedNotification(actorID, targetUserID string) error {
	return nil
}

func (m *MockNotificationService) MarkAsRead(notificationID int, userID string) error {
	return nil
}

func (m *MockNotificationService) MarkFollowRequestNotificationAsRead(actorID, targetUserID string) error {
	return nil
}

func (m *MockNotificationService) MarkGroupInviteNotificationAsRead(groupID, invitationID int, userID string) error {
	return nil
}

func TestFollowUser_InvalidMethod(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	mockNotificationService := &MockNotificationService{}
	followerHandler := NewFollowerHandler(db, mockNotificationService)

	req := httptest.NewRequest(http.MethodGet, "/api/users/follow", nil)
	// Inject mock user into context
	mockUser := &model.User{ID: "test-user"}
	ctx := ctxpkg.WithUser(req.Context(), mockUser)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	followerHandler.FollowUser(recorder, req)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", recorder.Code)
	}
}

func TestAcceptFollowRequest_InvalidMethod(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	mockNotificationService := &MockNotificationService{}
	followerHandler := NewFollowerHandler(db, mockNotificationService)

	req := httptest.NewRequest(http.MethodGet, "/api/follow/accept", nil)
	mockUser := &model.User{ID: "test-user"}
	ctx := ctxpkg.WithUser(req.Context(), mockUser)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	followerHandler.AcceptFollowRequest(recorder, req)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", recorder.Code)
	}
}

func TestDeclineFollowRequest_InvalidMethod(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	mockNotificationService := &MockNotificationService{}
	followerHandler := NewFollowerHandler(db, mockNotificationService)

	req := httptest.NewRequest(http.MethodGet, "/api/follow/decline", nil)
	mockUser := &model.User{ID: "test-user"}
	ctx := ctxpkg.WithUser(req.Context(), mockUser)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	followerHandler.DeclineFollowRequest(recorder, req)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", recorder.Code)
	}
}

func TestCancelFollowRequest_InvalidMethod(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	mockNotificationService := &MockNotificationService{}
	followerHandler := NewFollowerHandler(db, mockNotificationService)

	req := httptest.NewRequest(http.MethodGet, "/api/follow/cancel", nil)
	mockUser := &model.User{ID: "test-user"}
	ctx := ctxpkg.WithUser(req.Context(), mockUser)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	followerHandler.CancelFollowRequest(recorder, req)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", recorder.Code)
	}
}

func TestFollowUser_InvalidBody(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	mockNotificationService := &MockNotificationService{}
	followerHandler := NewFollowerHandler(db, mockNotificationService)

	req := httptest.NewRequest(http.MethodPost, "/api/users/follow", bytes.NewBuffer([]byte("notjson")))
	mockUser := &model.User{ID: "test-user"}
	ctx := ctxpkg.WithUser(req.Context(), mockUser)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	followerHandler.FollowUser(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", recorder.Code)
	}
}
