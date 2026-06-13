package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("A user with that email already exists")
	ErrDuplicateUsername = errors.New("A user with that username already exists")
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"` //we will not return password to the user
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
	Role_ID int64 `json:"role_id"`
	Role Role `json:"role"`
}

type password struct {
	text *string //plain text that user types
	hash []byte  //hashes and is stored in database
}

// function that receives password and generates hash
func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), 13)

	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash

	return nil

}

// function to comapare user password
func (p *password) ComparePassword(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
	// returns nil if match
	// returns error if no match
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `INSERT INTO users (username, email, password, role_id)
	VALUES($1, $2, $3, $4) RETURNING id, created_at
	`
	//creating a time out
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password.hash,
		user.Role_ID,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_email_key":
				return ErrDuplicateEmail
			case "users_username_key":
				return ErrDuplicateUsername
			}
		}
		return err
	}

	return nil
}

// Get user by id
func (s *UsersStore) GetById(ctx context.Context, userID int64) (*User, error) {
	query := ` SELECT u.id, u.username, u.email, u.created_at, r.* FROM users as u 
	JOIN 
	roles as r 
	ON (u.role_id = r.id)
	WHERE u.id = $1 AND u.is_Active=true`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	var user User
	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (s *UsersStore) CreateAndInvite(ctx context.Context, user *User, token string, expiry time.Duration) error {

	// withTx wraps everything in a transaction
	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		//transaction wrapper
		//-> Create the user
		//-> Create user user invite (if it fails, we roll back)

		//Step-1 - Create user in users table
		if err := s.Create(ctx, tx, user); err != nil {
			//if error roll back
			return err
		}

		//Step-2 create invitation in user_invitations table
		if err := s.createUserInvitation(ctx, tx, token, user.ID, expiry); err != nil {
			return err
		}

		return nil

	})
}

func (s *UsersStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, userID int64, expiry time.Duration) error {

	query := `insert into user_invitations (token,user_id,expiry) values ($1,$2,$3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(
		ctx,
		query,
		token,
		userID,
		time.Now().Add(expiry),
	)

	if err != nil {
		return err
	}

	return nil
}

// Activate user
func (s *UsersStore) ActivateUser(ctx context.Context, token string) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		//step-1 Find the user belonging to the token
		user, err := s.getUserFromInvitation(ctx, tx, token)

		if err != nil {
			return err
		}

		//step-2 update the user in the database as active
		user.IsActive = true
		if err := s.Update(ctx, tx, user); err != nil {
			return err

		}
		//step-3 Delete the user invitation token
		if err := s.DeleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UsersStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `select u.id,u.email,u.username,u.created_at,u.is_active 
	from user_invitations as i 
	inner join
	users as  u 
	 on u.id = i.user_id where (i.token= $1 AND expiry > $2);`

	//hash the token to compare
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}

	err := tx.QueryRowContext(
		ctx,
		query,
		hashToken,
		time.Now(),
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
		&user.IsActive,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UsersStore) Update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `update users set username=$1,email=$2,is_active =$3 where id =$4`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.IsActive,
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) DeleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `delete from user_invitations where user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *UsersStore) Delete(ctx context.Context, userID int64) error {
	query := `delete from users where id =$1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, username, password, created_at From users where email=$1 AND is_active=true`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password.hash,
		&user.CreatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, err
}
