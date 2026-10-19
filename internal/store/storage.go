package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	NotFoundError = errors.New("resource not found")
	ConflictError = errors.New("resource already exists")
)

type Users interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id int64) (*User, error)
}

type Posts interface {
	Create(ctx context.Context, posts *Post) error
	GetById(ctx context.Context, id int64) (*Post, error)
	Update(ctx context.Context, id int64, post Post) (*Post, error)
	DeleteById(ctx context.Context, id int64) error
	GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]*PostWithMetaData, error) 
}

type Comments interface {
	GetByPostId(ctx context.Context, postID int64) ([]Comment, error)
	Create(ctx context.Context, comment *Comment) error
}

type Followers interface {
	Follow(ctx context.Context, followerID, userId int64) error
	UnFollow(ctx context.Context, followerID, userId int64) error
}

type Storage struct {
	Users
	Posts
	Comments
	Followers
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStorage{db},
		Users:     &UserStorage{db},
		Comments:  &CommentStore{db},
		Followers: &FollowersStore{db},
	}
}
