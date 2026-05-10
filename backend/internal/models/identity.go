package models

import (
	"context"
)

//--------------------------------------------------------------------------------------|

type contextKey string
const UserKey contextKey = "user_identity"

//--------------------------------------------------------------------------------------|

type UserIdentity struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

//--------------------------------------------------------------------------------------|

func GetIdentity(ctx context.Context) *UserIdentity {
	identity, _ := ctx.Value(UserKey).(*UserIdentity)
	return identity
}

//--------------------------------------------------------------------------------------|

func WithIdentity(ctx context.Context, identity *UserIdentity) context.Context {
	return context.WithValue(ctx, UserKey, identity)
}