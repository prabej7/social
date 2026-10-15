package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id,omitempty"`
	UserName  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"-"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $2) RETURNING id, created_at
	`

	err := s.db.QueryRowContext(ctx, query,
		user.UserName,
		user.Email,
		user.Password,
	).Scan(user.ID, user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
