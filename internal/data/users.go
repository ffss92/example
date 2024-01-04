package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ffss92/example/internal/auth"
)

func (s Store) InsertUser(user *auth.User) error {
	query := `INSERT INTO users (username, password_hash)
	VALUES ($1, $2)
	RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, user.Username, user.PasswordHash).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users.username"):
			return auth.ErrDuplicateUsername
		default:
			return fmt.Errorf("failed to insert user to db: %w", err)
		}
	}

	return nil
}

func (s Store) GetUser(id int64) (*auth.User, error) {
	query := `
	SELECT id, username, password_hash, created_at, updated_at 
	FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user auth.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, auth.ErrNotFound
		default:
			return nil, fmt.Errorf("failed to get user from db: %w", err)
		}
	}

	return &user, nil
}

func (s Store) GetUserByEmail(email string) (*auth.User, error) {
	query := `
	SELECT id, username, password_hash, created_at, updated_at 
	FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user auth.User
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, auth.ErrNotFound
		default:
			return nil, fmt.Errorf("failed to get user from db: %w", err)
		}
	}

	return &user, nil
}

func (s Store) GetUserByUsername(username string) (*auth.User, error) {
	query := `
	SELECT id, username, password_hash, created_at, updated_at 
	FROM users WHERE username = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user auth.User
	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, auth.ErrNotFound
		default:
			return nil, fmt.Errorf("failed to get user from db: %w", err)
		}
	}

	return &user, nil
}

func (s Store) GetUserForToken(hash []byte, scope auth.Scope) (*auth.User, error) {
	query := `
	SELECT u.id, u.username, u.password_hash, u.created_at, u.updated_at 
	FROM tokens t
	INNER JOIN users u
	ON u.id = t.user_id
	WHERE t.hash = $1
	AND t.scope = $2
	AND t.expiry > CURRENT_TIMESTAMP`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user auth.User
	err := s.db.QueryRowContext(ctx, query, hash, scope).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, auth.ErrNotFound
		default:
			return nil, fmt.Errorf("failed to get user from db: %w", err)
		}
	}

	return &user, nil
}

func (s Store) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user from db: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected from result: %w", err)
	}
	if rows == 0 {
		return auth.ErrNotFound
	}

	return nil
}
