package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, string) (*Post, error)
		Delete(context.Context, string) error
		Update(context.Context, *Post, string) error
		GetUserFeed(context.Context, string, PaginatedFeedQuery) (*[]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetById(context.Context, string) (*User, error)
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostId(context.Context, string) (*[]Comment, error)
	}
	Follower interface {
		Follow(context.Context, string, string) error
		UnFollow(context.Context, string, string) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentStore{db},
		Follower: &FollowerStore{db},
	}

}
