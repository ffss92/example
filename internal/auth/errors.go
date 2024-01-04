package auth

import "errors"

var (
	ErrDuplicateUser      = errors.New("user already exists")
	ErrNotFound           = errors.New("user not found")
	ErrInvalidToken       = errors.New("token is invalid or expired")
	ErrInvalidCredentials = errors.New("invalid credentails")
)
