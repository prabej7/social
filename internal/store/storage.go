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
	Update(ctx context.Context, id int64, post Post) (*Post, error)
	DeleteById(ctx context.Context, id int64) error
}

type Comments interface {
	GetByPostId(ctx context.Context, postID int64) ([]Comment, error)
	Create(ctx context.Context, comment *Comment) error
}

type Storage struct {
	Users
	Posts
	Comments
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStorage{db},
		Users:    &UserStorage{db},
		Comments: &CommentStore{db},
	}
}
