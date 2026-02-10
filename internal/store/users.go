package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	Email     string `json:"email"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO User(Username,Password,Email) VALUES (?,?.?) RETURNING id, created_at`

	err := s.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
