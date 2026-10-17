package store

import (
	"context"
	"database/sql"
	"errors"
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
	VALUES ($1, $2, $3) RETURNING id, created_at
	`

	err := s.db.QueryRowContext(ctx, query,
		user.UserName,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) GetById(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, email
		FROM users
		WHERE id = $1
	`

	var user User
	if err := s.db.QueryRowContext(ctx, query, id).
		Scan(
			&user.ID,
			&user.UserName,
			&user.Email,
		); err != nil {
		switch {
		case errors.Is(sql.ErrNoRows, err):
			return nil, NotFoundError
		default:
			return nil, err
		}
	}

	return &user, nil

}
