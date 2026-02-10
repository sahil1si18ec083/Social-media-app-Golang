package store

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Post struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	fmt.Println("checking")
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
