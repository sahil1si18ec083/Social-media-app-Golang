package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

var ErrAlreadyFollowing = errors.New("already following user")
var ErrSelfFollow = errors.New("cannot follow yourself")
var ErrFollowNotFound = errors.New("follow relationship not found")

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID string, userID string) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `
		INSERT INTO followers (follower_id, user_id)
		VALUES ($1,$2)
	`

	_, err := s.db.ExecContext(ctx, query, followerID, userID)
	fmt.Println(err)
	if err != nil {

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {

			switch pqErr.Code {

			case "23505":
				return ErrAlreadyFollowing

			case "23514":
				if pqErr.Constraint == "no_self_follow" {
					return ErrSelfFollow
				}
			}
		}
		return err
	}
	return err

}

func (s *FollowerStore) UnFollow(ctx context.Context, followerID string, userID string) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `
		DELETE FROM  followers WHERE follower_id=$1 and user_id=$2
	`
	result, err := s.db.ExecContext(ctx, query, followerID, userID)

	fmt.Println(err)
	if err != nil {

		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrFollowNotFound
	}
	return nil

}
