package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	User      User   `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetByPostId(ctx context.Context, postID int64) ([]Comment, error) {
	var comments []Comment

	query := `
	SELECT comments.id, comments.post_id, comments.user_id, users.username, comments.created_at, comments.updated_at, comments.content
	FROM comments
    	JOIN users on comments.user_id = users.id
	WHERE
		post_id = $1
	ORDER BY comments.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, postID)

	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		c := Comment{
			User: User{},
		}

		if err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.User.UserName,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.Content,
		); err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments
		(content, post_id, user_id)
		VALUES
		($1, $2, $3)
		RETURNING id
	`

	if err := s.db.QueryRowContext(
		ctx,
		query,
		comment.Content,
		comment.PostID,
		comment.UserID,
	).Scan(&comment.ID); err != nil {
		return err
	}

	return nil

}
