package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

type User struct {
	ID        int64  `json:"id,omitempty"`
	UserName  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"-"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3) RETURNING id, created_at
	`

	err := tx.QueryRowContext(ctx, query,
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

func (s *UserStorage) CreateAndInvite(
	ctx context.Context,
	user *User,
	token string,
	expiresAt time.Duration,
) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, tx, token, expiresAt, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStorage) createUserInvitation(
	ctx context.Context,
	tx *sql.Tx,
	token string,
	exp time.Duration,
	userID int64,
) error {
	query := `
		INSERT INTO user_invitation (token, expires_at, user_id)
		VALUES ($1, $2, $3)
	`

	_, err := tx.ExecContext(ctx, query, token, time.Now().Add(exp), userID)

	if err != nil {
		return err
	}

	return nil

}

func (s *UserStorage) getUserByToken(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitation ui 
		ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expires_at > NOW()
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	var user User
	if err := tx.QueryRowContext(ctx, query, hashToken).
		Scan(
			&user.ID,
			&user.UserName,
			&user.Email,
			&user.CreatedAt,
			&user.IsActive,
		); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, NotFoundError
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserStorage) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE users
		SET username =$1, email = $2, is_active = $3
		WHERE id = $4
	`

	_, err := tx.ExecContext(ctx, query, user.UserName, user.Email, user.IsActive, user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) deleteInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
		DELETE FROM user_invitation
		WHERE user_id = $1
	`

	_, err := tx.ExecContext(ctx, query, userID)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserByToken(ctx, tx, token)

		if err != nil {
			return err
		}

		//update the user
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		//clean invitation
		if err := s.deleteInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil

	})

}
