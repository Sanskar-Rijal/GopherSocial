package db

import (
	"context"
	"database/sql"
	"time"
)


func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error){

	//connecting database
	db, err := sql.Open("postgres",addr)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	
	duration,err := time.ParseDuration(maxIdleTime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)
	
	//check if the connection is alive
	//if it takes more than 5 seconds, we will consider it as a failure
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err !=nil {
		return nil, err
	}

	//if everything is okay we will return database driver 
	return db, nil 
}