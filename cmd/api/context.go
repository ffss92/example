package main

import (
	"context"

	"github.com/ffss92/example/internal/auth"
)

type ctxKey string

const (
	userKey ctxKey = "user"
)

// Sets a user for a given context.
func (a api) setUser(ctx context.Context, user *auth.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// Gets a user from a given context. Panics if a user is not present.
func (a api) mustGetUser(ctx context.Context) *auth.User {
	user, ok := ctx.Value(userKey).(*auth.User)
	if ok {
		return user
	}

	panic("user not present in context")
}
