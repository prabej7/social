package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content,omitempty"`
	Title     string    `json:"title,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
	UserID    int64     `json:"user_id,omitempty"`
	CreatedAt string    `json:"createdAt,omitempty"`
	UpdatedAt string    `json:"updatedAt,omitempty"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments,omitempty"`
}

type PostStorage struct {
	db *sql.DB
}

func (s *PostStorage) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStorage) GetById(ctx context.Context, id int64) (*Post, error) {
	var post Post
	query := `
	SELECT id, title, content, created_at, updated_at, user_Id, tags, version
	FROM posts 
	WHERE id=$1
	`

	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, NotFoundError
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStorage) Update(ctx context.Context, id int64, post Post) (*Post, error) {
	query := `
		UPDATE posts
		SET title = $2, content = $3, tags = $4, version = version + 1
		WHERE id = $1 AND version=$5
		RETURNING id, title, content, tags, version
	`

	row := s.db.QueryRowContext(
		ctx,
		query,
		id,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.Version,
	)

	var updatedPost Post

	err := row.Scan(
		&updatedPost.ID,
		&updatedPost.Title,
		&updatedPost.Content,
		pq.Array(&updatedPost.Tags),
		&updatedPost.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundError
		}
		return nil, err
	}

	return &updatedPost, nil
}

func (s *PostStorage) DeleteById(ctx context.Context, id int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1
	`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return NotFoundError
	}

	return nil
}
