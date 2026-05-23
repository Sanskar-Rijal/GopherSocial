package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// For public we use capital letters
// so when we send the data to frontend it looks ugly so converting it into lower case
// using `json:"id"`
type Post struct {
	ID        int64     `json:"id"` // postgres generates this — you need it back
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserId    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"` // postgres generates this — you need it back
	UpdatedAt string    `json:"updated_at"` // postgres generates this — you need it back
	Comments  []Comment `json:"comments,omitempty"`
	User      *Cuser    `json:"user,omitempty"`
	Version   int       `json:"version"`
}

type PostWithMetaData struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostStore struct {
	db *sql.DB
}

// Create post
func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, user_id, tags)
	VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	//creating a time out
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,        // the SQL string with $1 $2 $3 $4 placeholders
		post.Content, // fills $1
		post.Title,   // fills $2
		post.UserId,  // fills $3
		//our tag is a slice of string, so we have to convert it into a format that postgres understands
		pq.Array(post.Tags), // fills $4
	).Scan(
		&post.ID,        // postgres generated id    → stored into post.ID
		&post.CreatedAt, // postgres generated time  → stored into post.CreatedAt
		&post.UpdatedAt, // postgres generated time  → stored into post.UpdatedAt
	)
	//After INSERT — postgres sends back: RETURNING id, created_at, updated_at
	//.Scan() receives those returned values and puts them INTO our post struct

	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetById(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at, version FROM posts WHERE id = $1`

	//creating a time out
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
		&post.Version,
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

// Delete post
func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	//creating a time out
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

// Update post using patch request
func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts SET content= $1, title= $2,  
	updated_at = NOW(), version = version +1 WHERE id = $3 AND version = $4 
	RETURNING updated_at, version`

	//creating a time out
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.ID,
		post.Version,
	).Scan(
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// two possibilities:
			// 1. post does not exist
			// 2. version mismatch - someone else updated first
			return ErrConflict
		default:
			return err
		}
	}

	return nil

}

// User feed
func (s *PostStore) GetUserFeed(ctx context.Context, userId int64, paginatedQuery *PaginatedQuery) ([]PostWithMetaData, error) {

	query :=
		`select 
		p.id,u.id as user_id,u.username, p.content,
		 p.title, p.tags, p.created_at, p.version, 
	count(c.id) as comment_count
	from posts as p 
	left join comments as c 
	on 
	p.id = c.post_id
	inner join users as u
	on 
	u.id = p.user_id
	inner join followers as f
	on 
	f.user_id = p.user_id 
	where (f.follower_id = $1  OR p.user_id = $1) AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') 
	AND (p.tags @> $5 OR $5 = '{}')
	group by p.id, u.id
	order by p.created_at ` + paginatedQuery.Sort + `
	LIMIT $2 OFFSET $3
	;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userId,
		paginatedQuery.Limit,
		paginatedQuery.Offset,
		paginatedQuery.Search,
		pq.Array(paginatedQuery.Tags),
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetaData
	var user Cuser

	for rows.Next() {
		var p PostWithMetaData

		err := rows.Scan(
			&p.ID,
			&p.UserId,
			&user.Username,
			&p.Content,
			&p.Title,
			pq.Array(&p.Tags),
			&p.CreatedAt,
			&p.Version,
			&p.CommentCount,
		)

		p.User = &user

		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}
	return feed, nil
}
