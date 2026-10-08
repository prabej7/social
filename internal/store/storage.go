package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	NotFoundError = errors.New("resource not found")
)

type Users interface {
	Create(ctx context.Context, user *User) error
}

type Posts interface {
	Create(ctx context.Context, posts *Post) error
	GetById(ctx context.Context, id int64) (*Post, error)
}

type Storage struct {
	Users
	Posts
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStorage{db},
		Users: &UsersStorage{db},
	}
}
