package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserId int64 `json:"user_id"` //my user id
	FollowerId int64 `json:"follower_id"` //id of the user who id following me
}

type FollowerStore struct {
	db *sql.DB
}


func (s *FollowerStore) FollowUser(ctx context.Context, followerId int64, userId int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	 
	_, err := s.db.ExecContext(
		ctx,
		query, 
		userId, 
		followerId,
	)
	if err != nil {
		//do something
		return err
	}
	return nil
}

func (s *FollowerStore) UnFollowUser(ctx context.Context, followerId int64, userId int64) error {
	query := `DELETE FROM followers where user_id = $1 AND follower_id = $2`
	
	ctx, cancel := context.WithTimeout(ctx,QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		userId,
		followerId,
	)
	if err != nil {
		return err
	}
	return nil
}