package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("Resource not found")
	ErrConflict          = errors.New("Resource Conflict")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64, *PaginatedQuery) ([]PostWithMetaData, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetById(context.Context, int64) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		ActivateUser(context.Context, string) error
	}

	Comments interface {
		GetCommentsFromPost(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
		Update(context.Context, *Comment) error
		Delete(context.Context, int64) error
	}
	Followers interface {
		FollowUser(context.Context, int64, int64) error
		UnFollowUser(context.Context, int64, int64) error
		GetFollowers(context.Context, int64) ([]Cuser, error)
		GetFollowing(context.Context, int64) ([]Cuser, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db: db},
		Users:     &UsersStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	//starting a transaction
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		//undo everything if anything fails
		_ = tx.Rollback()
		return err
	}

	//SAVE everything if all succeeded, make it permanent
	return tx.Commit()
}
