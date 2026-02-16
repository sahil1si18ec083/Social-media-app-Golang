package store

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Comment struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	UserID    int64  `json:"user_id"`
	PostID    int64  `json:"post_id"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, Comment *Comment) error {

	return nil
}
func (s *CommentStore) GetByPostId(ctx context.Context, postId string) (*[]Comment, error) {
	query := `select c.id, c.content,u.id,u.username from comments as c left join users as u on c.user_id = u.id
where c.post_id =$1  order by  c.updated_at desc`
	var comments []Comment
	rows, err := s.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var c Comment
		var u User
		err := rows.Scan(&c.ID,
			&c.Content,
			&u.ID,
			&u.Username)
		if err != nil {
			return nil, err
		}
		c.User = u

		c.UserID = u.ID
		int64val, err := strconv.ParseInt(postId, 10, 64)
		fmt.Println(int64val, "xx")
		if err != nil {
			return nil, err
		}
		c.PostID = int64val
		comments = append(comments, c)

	}

	return &comments, nil

}
