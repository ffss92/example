package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ffss92/example/internal/auth"
)

// Recovers from panics that happens in http handlers.
func (a api) recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Extracts the user from the Authorization header and adds to the request context.
func (a api) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Header not set.
		// Current user set to Anonymous and calls next.
		if authHeader == "" {
			ctx := a.setUser(r.Context(), auth.AnonymousUser)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Auth header should have 'Bearer <token>' format.
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.invalidTokenError(w, r)
			return
		}

		token := parts[1]
		// Tries to fetch the current user for the given token.
		user, err := a.auth.GetUserForToken(token, auth.ScopeAuthentication)
		if err != nil {
			switch {
			case errors.Is(err, auth.ErrInvalidToken):
				a.invalidTokenError(w, r)
			default:
				a.serverError(w, r, err)
			}
			return
		}

		// Add current user to the context and call the next handler.
		ctx := a.setUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Wraps an http handler and ensures that the current user is not anonymous.
func (a api) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := a.mustGetUser(r.Context())
		if user.IsAnonymous() {
			// Returns a 401 response.
			a.authRequiredError(w, r)
			return
		}

		// Calls the next handler.
		next.ServeHTTP(w, r)
	}
}
