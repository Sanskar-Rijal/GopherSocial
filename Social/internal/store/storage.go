package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("Resource not found")
	ErrConflict = errors.New("Resource Conflict")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error 
		Update(context.Context, *Post) error

	}
	Users interface {
		Create(context.Context, *User) error
		GetById(context.Context, int64) (*User, error)
	}

	Comments interface {
		GetCommentsFromPost(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface{
		FollowUser(context.Context, int64, int64) error
		UnFollowUser(context.Context, int64, int64) error
	}
}


func NewPostgresStorage(db *sql.DB) Storage{
	return Storage{
		Posts: &PostStore{db: db} ,
		Users: &UsersStore{db: db},
		Comments: &CommentStore{db:db},
		Followers: &FollowerStore{db:db},
	}
}