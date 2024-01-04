package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

type Scope string

const (
	// We could add diferent scopes, like password-reset, etc.
	ScopeAuthentication Scope = "authentication"
)

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	Scope     Scope     `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

// Creates a new token with a defined scope for a given user.
func (s Service) newToken(userId int64, ttl time.Duration, scope Scope) (*Token, error) {
	token := &Token{
		UserID: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		return nil, err
	}

	plainText := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes)

	token.PlainText = plainText
	token.Hash = s.hashToken(plainText)

	// Insert token to DB
	if err := s.storer.InsertToken(token); err != nil {
		return nil, err
	}

	return token, nil
}

// Creates a hash from a plain text token.
func (s Service) hashToken(plainText string) []byte {
	hash := sha256.Sum256([]byte(plainText))
	return hash[:]
}
