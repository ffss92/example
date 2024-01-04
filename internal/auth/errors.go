package auth

import "errors"

var (
	// TODO: Use fmt.Errorf to wrap a single duplicate error and add context.
	ErrDuplicateEmail     = errors.New("user email is already taken")
	ErrDuplicateUsername  = errors.New("user username is already taken")
	ErrNotFound           = errors.New("user not found")
	ErrInvalidToken       = errors.New("token is invalid or expired")
	ErrInvalidCredentials = errors.New("invalid credentails")
)
