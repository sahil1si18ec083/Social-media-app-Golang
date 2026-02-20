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
	UpdatedAt string `json:"updated_at"`
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

func (s *UsersStore) GetById(ctx context.Context, userId string) (*User, error) {
	query := `SELECT email, username,id, created_at, updated_at, password_hash from Users where id =$1`
	var user User
	err := s.db.QueryRowContext(ctx, query, userId).Scan(
		&user.Email,
		&user.Username,
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
