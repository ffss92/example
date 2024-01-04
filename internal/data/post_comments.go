package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ffss92/example/internal/posts"
)

func (s Store) InsertPostComment(comment *posts.Comment) error {
	query := `
	INSERT INTO post_comments (comment, user_id, post_id)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, comment.Comment, comment.UserID, comment.PostID).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert comment to db: %w", err)
	}

	return nil
}

func (s Store) GetComment(postId, commentId int64) (*posts.Comment, error) {
	query := `
	SELECT 
		pc.id, 
		pc.comment, 
		pc.user_id, 
		pc.post_id, 
		pc.created_at, 
		pc.updated_at,
		u.username
	FROM post_comments pc
	INNER JOIN users u
	ON u.id = pc.user_id
	WHERE pc.id = $1
	AND pc.post_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var comment posts.Comment
	err := s.db.QueryRowContext(ctx, query, commentId, postId).Scan(
		&comment.ID,
		&comment.Comment,
		&comment.UserID,
		&comment.PostID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Author,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, posts.ErrCommentNotFound
		default:
			return nil, fmt.Errorf("failed to get comment from db: %w", err)
		}
	}

	return &comment, nil
}

func (s Store) ListPostComments(postId int64, limit, offset int) ([]*posts.Comment, error) {
	query := `
	SELECT 
		pc.id, 
		pc.comment, 
		pc.user_id, 
		pc.post_id, 
		pc.created_at, 
		pc.updated_at,
		u.username
	FROM post_comments pc
	INNER JOIN users u
	ON u.id = pc.user_id
	WHERE pc.post_id = $1
	LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var comments []*posts.Comment

	rows, err := s.db.QueryContext(ctx, query, postId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment posts.Comment

		err := rows.Scan(
			&comment.ID,
			&comment.Comment,
			&comment.UserID,
			&comment.PostID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Author,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan post comment: %w", err)
		}

		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list comments from db: %w", err)
	}

	return comments, nil
}

func (s Store) DeletePostComment(commentId int64) error {
	query := `DELETE FROM post_comments WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, commentId)
	if err != nil {
		return fmt.Errorf("failed to delete comment from db: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return posts.ErrCommentNotFound
	}

	return nil
}

func (s Store) UpdatePostComment(comment *posts.Comment) error {
	query := `
	UPDATE post_comments
	SET
		comment = $1,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $2
	RETURNING updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, comment.Comment, comment.ID).Scan(
		&comment.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment to db: %w", err)
	}

	return nil
}

func (s Store) CountPostComments(postId int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM post_comments WHERE post_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, postId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get post count from db: %w", err)
	}

	return count, nil
}
