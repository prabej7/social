package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	UserID    int64     `json:"user_id"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
	Comments  []Comment `json:"comments"`
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
	SELECT id, title, content, created_at, updated_at, user_Id, tags 
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
