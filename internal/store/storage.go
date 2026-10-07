package store

import (
	"context"
	"database/sql"
)

type Users interface {
	Create(ctx context.Context, user *User) error
}

type Posts interface {
	Create(ctx context.Context, posts *Post) error
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
