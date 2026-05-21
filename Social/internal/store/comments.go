package store

import (
	"context"
	"database/sql"
	"errors"
)

//creating struct to send user details

type Cuser struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
}


type Comment struct {
	ID int64 `json:"id"`
	PostID int64 `json:"post_id"`
	UserID int64 `json:"user_id"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	User *Cuser `json:"user,omitempty"`
}

type CommentStore struct {
	db *sql.DB
}


func (s *CommentStore) Create(ctx context.Context, comment *Comment) error{
	query := `INSERT INTO comments (post_id, user_id, content)
	 VALUES ($1,$2,$3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err !=nil {
		return err
	}
	return nil 
}


func (s *CommentStore) GetCommentsFromPost(ctx context.Context, postId int64) ([]Comment, error) {
	query := `select c.id,c.post_id,c.user_id,c.content,c.created_at,u.id,u.username from comments as c inner join 
	users as u on u.id = c.user_id where c.post_id = $1 order by c.created_at desc;`

	rows, err := s.db.QueryContext(ctx,query,postId)
	
	if err != nil {
		return nil, err
	}

	//we have to parse rows so we are closing them as well 
	defer rows.Close()

	comments := []Comment{}

	for rows.Next(){
		var c Comment
		var u Cuser

		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&u.ID,
			&u.Username,
		)

		if err != nil {
			return nil, err
		}
		c.User = &u
		comments = append(comments, c)
	}
	return comments, nil
}

func (s *CommentStore) Delete(ctx context.Context, commentId int64) error {
	query := `DELETE FROM comments WHERE id =$1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(
		ctx,
		query,
		commentId,
	)

	if err != nil {
		return err
	}

	//checking if rows are affected or not, if none of the rows are affected,
	//that means nothing is deleteed, that means comment doesn't exists
	rows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *CommentStore) Update(ctx context.Context, comment *Comment) error {
	query := `UPDATE comments SET content = $1 WHERE id = $2 RETURNING created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.Content,
		comment.ID,
	).Scan(
		&comment.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default: 
			return err
		}
	}
	return nil
}


