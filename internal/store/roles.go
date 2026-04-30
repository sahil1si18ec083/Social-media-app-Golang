package store

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int64  `json:"level"`
	CreatedAt   string `json:"created_at"`
}

type RoleStore struct {
	db *sql.DB
}

func (s *RoleStore) GetByRolename(ctx context.Context, rolename string) (*Role, error) {
	var role Role
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `Select  id, name, description, level,  created_at FROM roles where name = $1`
	err := s.db.QueryRowContext(ctx, query, rolename).Scan(&role.ID, &role.Name, &role.Description, &role.Level, &role.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &role, nil

}
