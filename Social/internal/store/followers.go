package store

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
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
		//do something because when user click follow then again click follow then 
		//it will give error because primary key is uniuqe
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return ErrConflict
	}
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


func (s *FollowerStore) GetFollowers(ctx context.Context, userId int64) ([]Cuser, error){
	query := `SELECT u.id, u.username FROM users AS u INNER JOIN followers AS f ON 
	f.follower_id = u.id WHERE f.user_id =$1`;
	//getting all users that follow userID

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userId,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var followers []Cuser

	for rows.Next(){
		var u Cuser
		err := rows.Scan(
			&u.ID,
			&u.Username,
		)

		if err != nil {
			return nil, err
		}
		followers = append(followers, u)
	}
	return followers, nil
}

func (s *FollowerStore) GetFollowing(ctx context.Context, userId int64) ([]Cuser, error){
	query := `SELECT u.id, u.username FROM users AS u INNER JOIN followers AS f ON 
	f.user_id = u.id WHERE f.follower_id =$1`;
	
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []Cuser

	for rows.Next(){
		var u Cuser
		err := rows.Scan(
			&u.ID,
			&u.Username,
		)
		if err != nil {
			return nil, err
		}
		following = append(following, u)
	}
	return following, nil
}