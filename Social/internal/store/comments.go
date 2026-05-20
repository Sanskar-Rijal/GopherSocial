package store

import (
	"context"
	"database/sql"
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
	User Cuser `json:"user"`
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
	query := `select c.id,c.user_id,c.content,c.created_at,u.id,u.username from comments as c inner join 
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
		// c.User = User{}

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.User.ID,
			&c.User.Username,
		)

		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}