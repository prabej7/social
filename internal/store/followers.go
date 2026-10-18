package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowersStore struct {
	db *sql.DB
}

func (s *FollowersStore) Follow(ctx context.Context, followerId int64, userId int64) error {
	query := `
		INSERT INTO followers
		(follower_id, user_id)
		VALUES
		($1, $2)
	`

	_, err := s.db.ExecContext(ctx, query, followerId, userId)

	if err != nil {
		return err
	}

	return nil
}

func (s *FollowersStore) UnFollow(ctx context.Context, followerId int64, userId int64) error {
	query := `
		DELETE FROM followers
		WHERE follower_id = $1 AND user_id = $2
	`

	_, err := s.db.ExecContext(ctx, query, followerId, userId)

	if err != nil {
		return err
	}

	return nil
}
