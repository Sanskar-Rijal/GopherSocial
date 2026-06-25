package cache

import (
	"context"
	"social/internal/store"

	"github.com/redis/go-redis/v9"
)

type Storage struct{
	Users interface {
		GetUser(context.Context, int64) (*store.User, error)
		SetUser(context.Context, *store.User) error 
		DeleteUser(context.Context, int64)
	}
}


func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}