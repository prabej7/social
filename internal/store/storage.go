package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	NotFoundError = errors.New("resource not found")
	ConflictError = errors.New("resource already exists")
)

type Users interface {
	Create(ctx context.Context, tx *sql.Tx, user *User) error
	GetById(ctx context.Context, id int64) (*User, error)
	CreateAndInvite(
		ctx context.Context,
		user *User,
		token string,
		expiresAt time.Duration,
	) error
	Activate(ctx context.Context, token string) error
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

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
