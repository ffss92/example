package data

import (
	"context"
	"fmt"
	"time"

	"github.com/ffss92/example/internal/posts"
)

func (s Store) InsertPostLike(userId, postId int64) error {
	query := `INSERT INTO post_likes (user_id, post_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, postId)
	if err != nil {
		switch {
		case isUniqueError(err, "post_likes", "user_id", "post_id"):
			return posts.ErrAlreadyLiked
		default:
			return fmt.Errorf("failed to insert post like in db: %w", err)
		}
	}

	return nil
}

func (s Store) DeletePostLike(userId, postId int64) error {
	query := `DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, postId, userId)
	if err != nil {
		return fmt.Errorf("failed to insert post like in db: %w", err)
	}

	return nil
}
