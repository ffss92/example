package data

import (
	"context"
	"fmt"
	"time"

	"github.com/ffss92/example/internal/auth"
)

func (s Store) InsertToken(token *auth.Token) error {
	query := `
	INSERT INTO tokens (hash, scope, user_id, expiry)
	VALUES ($1, $2, $3, $4)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, token.Hash, token.Scope, token.UserID, token.Expiry)
	if err != nil {
		return fmt.Errorf("failed to insert token to db: %w", err)
	}

	return nil
}

func (s Store) DeleteTokensForUser(userId int64, scope auth.Scope) error {
	query := `
	DELETE FROM tokens
	WHERE user_id = $1
	AND scope = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, scope)
	if err != nil {
		return fmt.Errorf("failed to delete user tokens from db: %w", err)
	}

	return nil
}
