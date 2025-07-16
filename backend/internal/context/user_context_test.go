package context

import (
	"backend/internal/model"
	"context"
	"testing"
	"time"
)

func TestWithUser(t *testing.T) {
	ctx := context.Background()
	user := &model.User{
		ID:        "test-user-id",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Add user to context
	ctxWithUser := WithUser(ctx, user)

	// Verify user was added
	retrievedUser, ok := GetUser(ctxWithUser)
	if !ok {
		t.Fatal("Expected user to be found in context")
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, retrievedUser.ID)
	}

	if retrievedUser.Email != user.Email {
		t.Errorf("Expected user email %s, got %s", user.Email, retrievedUser.Email)
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()

	// Test with no user in context
	_, ok := GetUser(ctx)
	if ok {
		t.Error("Expected no user in empty context")
	}

	// Test with user in context
	user := &model.User{
		ID:    "test-user-id",
		Email: "test@example.com",
	}
	ctxWithUser := WithUser(ctx, user)

	retrievedUser, ok := GetUser(ctxWithUser)
	if !ok {
		t.Fatal("Expected user to be found in context")
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, retrievedUser.ID)
	}
}

func TestWithSessionID(t *testing.T) {
	ctx := context.Background()
	sessionID := "test-session-id"

	// Add session ID to context
	ctxWithSession := WithSessionID(ctx, sessionID)

	// Verify session ID was added
	retrievedSessionID, ok := GetSessionID(ctxWithSession)
	if !ok {
		t.Fatal("Expected session ID to be found in context")
	}

	if retrievedSessionID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, retrievedSessionID)
	}
}

func TestGetSessionID(t *testing.T) {
	ctx := context.Background()

	// Test with no session ID in context
	_, ok := GetSessionID(ctx)
	if ok {
		t.Error("Expected no session ID in empty context")
	}

	// Test with session ID in context
	sessionID := "test-session-id"
	ctxWithSession := WithSessionID(ctx, sessionID)

	retrievedSessionID, ok := GetSessionID(ctxWithSession)
	if !ok {
		t.Fatal("Expected session ID to be found in context")
	}

	if retrievedSessionID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, retrievedSessionID)
	}
}

func TestMustGetUser(t *testing.T) {
	// Test with user in context - should not panic
	user := &model.User{
		ID:        "test-user-id",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		DOB:       time.Now(),
		CreatedAt: time.Now(),
	}
	ctx := WithUser(context.Background(), user)

	retrievedUser := MustGetUser(ctx)
	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, retrievedUser.ID)
	}

	// Test with no user in context - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected MustGetUser to panic when no user in context")
		}
	}()
	MustGetUser(context.Background())
}

func TestMustGetSessionID(t *testing.T) {
	// Test with session ID in context - should not panic
	sessionID := "test-session-id"
	ctx := WithSessionID(context.Background(), sessionID)

	retrievedSessionID := MustGetSessionID(ctx)
	if retrievedSessionID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, retrievedSessionID)
	}

	// Test with no session ID in context - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected MustGetSessionID to panic when no session ID in context")
		}
	}()
	MustGetSessionID(context.Background())
}

func TestContextKeyCollisions(t *testing.T) {
	ctx := context.Background()

	// Add both user and session ID to context
	user := &model.User{
		ID:    "test-user-id",
		Email: "test@example.com",
	}
	sessionID := "test-session-id"

	ctx = WithUser(ctx, user)
	ctx = WithSessionID(ctx, sessionID)

	// Verify both can be retrieved independently
	retrievedUser, userOk := GetUser(ctx)
	retrievedSessionID, sessionOk := GetSessionID(ctx)

	if !userOk {
		t.Error("Expected user to be found in context")
	}
	if !sessionOk {
		t.Error("Expected session ID to be found in context")
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, retrievedUser.ID)
	}
	if retrievedSessionID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, retrievedSessionID)
	}
}
