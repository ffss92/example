package auth

import "github.com/go-playground/validator/v10"

type Storer interface {
	// Users
	InsertUser(*User) error
	GetUser(int64) (*User, error)
	GetUserByEmail(string) (*User, error)
	GetUserByUsername(string) (*User, error)
	GetUserForToken([]byte, Scope) (*User, error)
	DeleteUser(int64) error

	// Tokens
	DeleteTokensForUser(int64, Scope) error
	InsertToken(*Token) error
}

type Service struct {
	storer   Storer
	validate *validator.Validate
}

func NewService(storer Storer, validate *validator.Validate) Service {
	return Service{
		storer:   storer,
		validate: validate,
	}
}
