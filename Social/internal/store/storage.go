package store

import (
	"context"
	"database/sql"
	"errors"
)

var (ErrNotFound = errors.New("Record not found"))

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)

	}
	Users interface {
		Create(context.Context, *User) error
	}

	Comments interface {
		GetCommentsFromPost(context.Context, int64) ([]Comment, error)
	}
}


func NewPostgresStorage(db *sql.DB) Storage{
	return Storage{
		Posts: &PostStore{db: db} ,
		Users: &UsersStore{db: db},
		Comments: &CommentStore{db:db},
	}
}