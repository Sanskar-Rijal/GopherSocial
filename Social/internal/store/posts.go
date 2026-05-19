package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

//For public we use capital letters
//so when we send the data to frontend it looks ugly so converting it into lower case
//using `json:"id"`
type Post struct {
	ID int64 `json:"id"` // postgres generates this — you need it back
	Content string `json:"content"`
	Title string `json:"title"`
	UserId int64 `json:"user_id"`
	Tags []string `json:"tags"`
	CreatedAt string `json:"created_at"` // postgres generates this — you need it back
	UpdatedAt string `json:"updated_at"` // postgres generates this — you need it back
	Comments []Comment `json:"comments"`
}

type PostStore struct {
	db *sql.DB
}

//Create post
func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, user_id, tags)
	VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRowContext(
		ctx,
		 query, // the SQL string with $1 $2 $3 $4 placeholders
		post.Content,    // fills $1
		post.Title,   // fills $2
		post.UserId,   // fills $3
		//our tag is a slice of string, so we have to convert it into a format that postgres understands
		pq.Array(post.Tags), // fills $4
		).Scan(
			&post.ID, // postgres generated id    → stored into post.ID
			&post.CreatedAt, // postgres generated time  → stored into post.CreatedAt
			&post.UpdatedAt, // postgres generated time  → stored into post.UpdatedAt
		)
		//After INSERT — postgres sends back: RETURNING id, created_at, updated_at
		//.Scan() receives those returned values and puts them INTO our post struct

		if err !=nil {
			return err
		}
		return  nil
}

func (s *PostStore) GetById(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at FROM posts WHERE id = $1`

	var post Post
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.Content,
		&post.Title, 
		&post.UserId,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrNotFound
			default: 
				return nil, err
			} 
	}
	return &post, nil
}

//Delete post 
func (s * PostStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, postID)

	if err != nil {
		return err
	}

	//Checking if any row is affected or not
	//if none of the rows is affected, that means nothing is deleted, it means post doesn't exists
	rows, err := res.RowsAffected()

	if err != nil {
		return err 
	}

	//if rows = 0 then nothing is delted, so we return "post not found"
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}


//Update post using patch request 
func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts SET content= $1, title= $2,  updated_at = NOW() WHERE id = $3 RETURNING updated_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.ID,
	).Scan(
		&post.UpdatedAt,
	)

	if err != nil {
		switch{
		case errors.Is(err,sql.ErrNoRows):
			return ErrNotFound //it means post doesn't exits
		default:
			return err
		}
	}

	return err;

}