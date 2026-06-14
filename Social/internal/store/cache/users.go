package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"social/internal/store"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rdb *redis.Client
}

const userExpTime = time.Minute * 5

func (s *UserStore) GetUser(ctx context.Context, userID int64) (*store.User, error){

	//{user:42} making it a key {user-42} i.e a key
	cacheKey := fmt.Sprintf("user-%d", userID)
	data ,err := s.rdb.Get(ctx, cacheKey).Result()

	if err == redis.Nil{
		return nil, nil 
	} else if err != nil {
		return nil, err 
	}

	//data is of type string so unmarshaling it to json format 
	var user store.User
	if data != ""{
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil 
}


func (s  *UserStore) SetUser(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%d", user.ID)

	json, err := json.Marshal(user)

	if err != nil {
		return err 
	}

	return s.rdb.Set(ctx, cacheKey,json, userExpTime).Err() //TTL (time to live is 5 minutes)
}