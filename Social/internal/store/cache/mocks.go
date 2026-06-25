package cache

import (
	"context"
	"social/internal/store"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {

}


func(m *MockUserStore) GetUser(ctx context.Context, userID int64) (*store.User, error){
	return nil, nil 
}
func(m *MockUserStore) SetUser(context.Context, *store.User) error {
	return nil 
}
func(m *MockUserStore) DeleteUser(context.Context, int64){

}