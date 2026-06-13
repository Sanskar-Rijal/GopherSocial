package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Level int `json:"level"`
	Description string `json:"description"`
}

type RoleStore struct{
	db *sql.DB
}


func (store *RoleStore) GetByName(ctx context.Context, slug string)(*Role, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `SELECT id, name, level, description FROM roles WHERE name =$1`

	role := &Role{}

	err := store.db.QueryRowContext(
		ctx,
		query,
		slug,
	).Scan(
		&role.ID, 
		&role.Name,
		&role.Level,
		&role.Description,
	)

	if err != nil {
		return nil, err 
	}

	return role, nil

}

