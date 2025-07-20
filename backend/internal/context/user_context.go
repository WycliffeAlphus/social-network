package context

import (
	"context"
	"backend/internal/model"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserContextKey is the key used to store user information in request context
	UserContextKey ContextKey = "user"
	// SessionIDContextKey is the key used to store session ID in request context
	SessionIDContextKey ContextKey = "session_id"
)

// WithUser adds a user to the context
func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// WithSessionID adds a session ID to the context
func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDContextKey, sessionID)
}

// MustGetUser retrieves the user from context and panics if not found
// This should only be used in handlers that are protected by auth middleware
func MustGetUser(ctx context.Context) *model.User {
	user, ok := ctx.Value(UserContextKey).(*model.User)
	if !ok {
		panic("user not found in context - ensure auth middleware is applied")
	}
	return user
}

// MustGetSessionID retrieves the session ID from context and panics if not found
// This should only be used in handlers that are protected by auth middleware
func MustGetSessionID(ctx context.Context) string {
	sessionID, ok := ctx.Value(SessionIDContextKey).(string)
	if !ok {
		panic("session ID not found in context - ensure auth middleware is applied")
	}
	return sessionID
}
