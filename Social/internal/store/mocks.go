package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {

	return  Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {

}

func (m *MockUserStore) Create(ctx context.Context,tx *sql.Tx, u *User)error {
	return nil
}

func(m *MockUserStore) Delete(ctx context.Context, ID int64) error{
	return nil 
}

func(m *MockUserStore) GetById(context.Context, int64) (*User, error){
	return &User{},nil
}
func(m *MockUserStore) CreateAndInvite(context.Context, *User, string, time.Duration) error{
	return nil 
}
func(m *MockUserStore) ActivateUser(context.Context, string) error{
	return nil 
}
func(m *MockUserStore) GetByEmail(context.Context, string) (*User, error) {
	return &User{}, nil 
}