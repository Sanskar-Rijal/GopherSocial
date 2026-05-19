package store

import (
	"context"
	"database/sql"
)


type Comment struct {
	ID int64 `json:"id"`
	PostID int64 `json:"post_id"`
	UserID int64 `json:"user_id"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	User User `json:"user"`
}

type CommentStore struct {
	db *sql.DB
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