package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Email     string   `json:"email"`
	IsActive  bool     `json:"activated"`
	Role      int64    `json:"role"`
}
type Password struct {
	Text *string
	Hash []byte
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User, tx *sql.Tx) error {
	query := `INSERT INTO users (username, password_hash, email, role_id) VALUES ($1, $2, $3,$4) RETURNING id, created_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := tx.QueryRowContext(ctx, query, user.Username, user.Password.Hash, user.Email, user.Role).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
func (p *Password) SetPassword(text string, user *User) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Text = &text
	p.Hash = hash
	return nil

}

func (s *UsersStore) GetById(ctx context.Context, userId string) (*User, error) {
	query := `SELECT email, username,id, created_at, updated_at, password_hash,role_id from Users where id =$1`
	var user User
	err := s.db.QueryRowContext(ctx, query, userId).Scan(
		&user.Email,
		&user.Username,
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Password.Hash,
		&user.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
func (s *UsersStore) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT email, username,id, created_at, updated_at, password_hash,role_id from Users where username =$1`
	var user User
	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.Email,
		&user.Username,
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Password.Hash,
		&user.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
func (s *UsersStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		if err := s.Create(ctx, user, tx); err != nil {
			fmt.Println(err)
			fmt.Println("yoo")
			return err
		}
		if err := s.CreateUserInvitation(ctx, user, tx, token, exp); err != nil {
			return err
		}
		return nil

	})
}

func (s *UsersStore) CreateUserInvitation(ctx context.Context, user *User, tx *sql.Tx, token string, exp time.Duration) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expiry)
		VALUES ($1, $2, NOW() + $3 * INTERVAL '1 second')
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, user.ID, int64(exp.Seconds()))
	if err != nil {
		return err
	}

	return nil
}
func (s *UsersStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username = $1, email = $2, activated = $3 WHERE id = $4`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}
func (s *UsersStore) Activate(ctx context.Context, token string) (*User, error) {

	var activatedUser *User

	err := withTx(s.db, ctx, func(tx *sql.Tx) error {
		// 1. find the user that this token belongs to
		user, err := s.getUserFromInvitation(ctx, tx, token)

		if err != nil {
			return err
		}

		// 2. update the user
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {

			return err
		}

		// 3. clean the invitations
		if err := s.deleteUserInvitations(ctx, tx, user.ID); err != nil {

			return err
		}
		activatedUser = user

		return nil
	})
	if err != nil {
		return nil, err
	}
	return activatedUser, nil

}

func (s *UsersStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {

	query := `SELECT u.id, u.username, u.email, u.created_at, u.activated from users as u
	join user_invitations as ui 
	ON ui.user_id = u.id
	where ui.token=$1 and 
	ui.expiry > $2
	
	;
	`
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
func (s *UsersStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userid int64) error {

	query := `DELETE FROM user_invitations WHERE user_id =$1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, userid)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil

}
