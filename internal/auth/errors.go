package auth

import "errors"

var (
	ErrDuplicateEmail     = errors.New("user email is already taken")
	ErrNotFound           = errors.New("user not found")
	ErrInvalidToken       = errors.New("token is invalid or expired")
	ErrInvalidCredentials = errors.New("invalid credentails")
)
