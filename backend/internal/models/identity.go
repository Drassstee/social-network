package models

import (
	"context"
)

//--------------------------------------------------------------------------------------|

// contextKey is an unexported type to prevent context key collisions.
type contextKey string

// UserKey is the context key for storing the authenticated user's identity.
const UserKey contextKey = "user_identity"

//--------------------------------------------------------------------------------------|

// UserIdentity holds the minimal set of claims for an authenticated user.
type UserIdentity struct {
	ID       int64  `json:"id"`
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
