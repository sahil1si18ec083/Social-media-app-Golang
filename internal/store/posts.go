package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
)

var ErrNotFound = errors.New("resource not found")
var QueryTimeoutDuration = time.Second * 1

type Post struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Comment   []Comment
	Version   int64 `json:"version"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	fmt.Println("checking")
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `INSERT INTO posts(title,content,user_id,tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.UserID, post.Tags).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	fmt.Println(err)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetById(ctx context.Context, postId string) (*Post, error) {
	var post Post
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `Select version, id, user_id, title, content, tags, created_at,  updated_at FROM posts where id = $1`
	err := s.db.QueryRowContext(ctx, query, postId).Scan(&post.Version, &post.ID, &post.UserID, &post.Title, &post.Content, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt)
	fmt.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &post, nil

}
func (s *PostStore) Delete(ctx context.Context, postId string) error {
	fmt.Println("checking")
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `DELETE FROM posts where id = $1 `
	res, err := s.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(rowsAffected)

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post, postid string) error {

	fmt.Println("checking")
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `
	UPDATE posts
	SET title = $1,
	    content = $2,
	    version = version + 1
	WHERE id = $3
	AND version = $4
	`

	result, err := s.db.ExecContext(
		ctx,
		query,
		post.Title,
		post.Content,
		postid,
		post.Version,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	fmt.Println(err)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil

}
