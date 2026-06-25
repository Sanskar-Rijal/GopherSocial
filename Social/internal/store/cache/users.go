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

type Sanskar struct{
	name string
}

const userExpTime = time.Minute * 5

//function to get user from the cache 
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


//function to store user in the cache
func (s  *UserStore) SetUser(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%d", user.ID)

	json, err := json.Marshal(user)

	if err != nil {
		return err 
	}

	return s.rdb.Set(ctx, cacheKey,json, userExpTime).Err() //TTL (time to live is 5 minutes)
}


//function to delete user from cache 
func (s *UserStore) DeleteUser(ctx context.Context, userID int64){
	cacheKey := fmt.Sprintf("user-%d",userID)
	s.rdb.Del(ctx, cacheKey)
}