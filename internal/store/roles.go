package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Level       int    `json:"level,omitempty"`
	Description string `json:"description,omitempty"`
}

type RoleStore struct {
	db *sql.DB
}

func (s *RoleStore) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `
		SELECT id, name, level, description
		FROM roles
		WHERE name = $1
	`

	var role Role
	if err := s.db.QueryRowContext(ctx, query, name).
		Scan(
			&role.ID,
			&role.Name,
			&role.Level,
			&role.Description,
		); err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *RoleStore) GetByID(ctx context.Context, id int64) (*Role, error) {
	query := `
		SELECT id, name, level, description
		WHERE id = $1
	`

	var role Role
	if err := s.db.QueryRowContext(ctx, query, id).
		Scan(
			&role.ID,
			&role.Name,
			&role.Level,
			&role.Description,
		); err != nil {
		return nil, err
	}

	return &role, nil
}
