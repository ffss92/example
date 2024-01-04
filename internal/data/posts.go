package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ffss92/example/internal/posts"
)

func (s Store) InsertPost(post *posts.Post) error {
	query := `
	INSERT INTO posts (title, content, user_id) 
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.UserID).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case isUniqueError(err, "posts", "title", "user_id"):
			return fmt.Errorf("%w: title", posts.ErrDuplicate)
		default:
			return fmt.Errorf("failed to insert post to db: %w", err)
		}
	}

	return nil
}

// TODO: Include like count
func (s Store) GetPost(id int64) (*posts.Post, error) {
	query := `
	SELECT 
		p.id, 
		p.title, 
		p.content, 
		p.user_id, 
		p.created_at, 
		p.updated_at, 
		u.username
	FROM posts p
	INNER JOIN users u
	ON u.id = p.user_id
	WHERE p.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var post posts.Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Author,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, posts.ErrNotFound
		default:
			return nil, fmt.Errorf("failed to get post from db: %w", err)
		}
	}

	return &post, nil
}

func (s Store) CountPosts() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM posts`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get post count from db: %w", err)
	}

	return count, nil
}

func (s Store) ListPosts(limit, offset int) ([]*posts.Post, error) {
	query := `
	SELECT 
		p.id, 
		p.title, 
		p.content, 
		p.user_id, 
		p.created_at, 
		p.updated_at, 
		u.username,
		COUNT(pl.post_id)
	FROM posts p
	INNER JOIN users u
	ON u.id = p.user_id
	LEFT JOIN post_likes pl
	ON pl.post_id = p.id
	GROUP BY p.id
	ORDER BY p.created_at DESC
	LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var p []*posts.Post
	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post posts.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserID,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Author,
			&post.Likes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row into post: %w", err)
		}

		p = append(p, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get posts from db: %w", err)
	}

	return p, nil
}

func (s Store) UpdatePost(post *posts.Post) error {
	query := `
	UPDATE posts 
	SET 
		title = $1, 
		content = $2, 
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	RETURNING updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID).Scan(
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case isUniqueError(err, "posts", "title", "user_id"):
			return fmt.Errorf("%w: title", posts.ErrDuplicate)
		default:
			return fmt.Errorf("failed to update posts in db: %w", err)
		}
	}

	return nil
}

func (s Store) DeletePost(id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post from db: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	} else if rows == 0 {
		return posts.ErrNotFound
	}

	return nil
}
