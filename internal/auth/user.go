package auth

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	authTokenTTL = time.Hour * 24 * 2 // 2 days.
)

var (
	AnonymousUser = &User{}
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s Service) CreateUser(params CreateUserParams) (*User, error) {
	pwHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to create password hash: %w", err)
	}

	user := &User{
		Email:        params.Email,
		PasswordHash: string(pwHash),
	}

	if err := s.storer.InsertUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s Service) Get(id int64) (*User, error) {
	user, err := s.storer.GetUser(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s Service) GetUserForToken(plainText string, scope Scope) (*User, error) {
	hash := s.hashToken(plainText)
	user, err := s.storer.GetUserForToken(hash, scope)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, ErrInvalidToken
		default:
			return nil, err
		}
	}

	return user, nil
}

type CredentialsParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s Service) Authenticate(creds CredentialsParam) (*Token, error) {
	user, err := s.storer.GetUserByEmail(creds.Email)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, ErrInvalidCredentials
		default:
			return nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return nil, ErrInvalidCredentials
		default:
			return nil, err
		}
	}

	token, err := s.newToken(user.ID, authTokenTTL, ScopeAuthentication)
	if err != nil {
		return nil, err
	}

	return token, nil
}
