package models

import (
	"context"
	"time"
)

//--------------------------------------------------------------------------------------|

// contextKey is an unexported type to prevent context key collisions.
type contextKey string

// UserKey is the context key for storing the authenticated user's identity.
const UserKey contextKey = "user_identity"

//--------------------------------------------------------------------------------------|

// UserIdentity holds the minimal set of claims for an authenticated user.
type UserIdentity struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

//--------------------------------------------------------------------------------------|

// GetIdentity extracts the UserIdentity from the request context.
// Returns nil if no identity is present (unauthenticated request).
func GetIdentity(ctx context.Context) *UserIdentity {
	identity, _ := ctx.Value(UserKey).(*UserIdentity)
	return identity
}

//--------------------------------------------------------------------------------------|

// WithIdentity returns a new context carrying the given UserIdentity.
func WithIdentity(ctx context.Context, identity *UserIdentity) context.Context {
	return context.WithValue(ctx, UserKey, identity)
}

//--------------------------------------------------------------------------------------|

// Session represents an active user session stored in the database.
type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
