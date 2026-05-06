package main

import (
	"log"
	"social/internal/db"
	"social/internal/env"
	"social/internal/store"
)

func main() {

	cfg := config{
			addr: env.GetString("ADDR",":8080"),
			db: dbConfig{
				addr: env.GetString("DB_ADDR","postgres://user:adminpassword@localhost/social?sslmode=disable"),
				maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS",30),
				maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS",30),
				maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME","15min"),
			},
		}

	// app := &application{
	// 	config: config{
	// 		addr: env.GetString("ADDR",":8080"),
	// 		db: dbConfig{
	// 			addr: env.GetString("DB_ADDR","postgres://user:adminpassword@localhost/social?sslmode=disable"),
	// 			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS",30),
	// 			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS",30),
	// 			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME","15min"),
	// 		},
	// 	},
	// 	store: store,
	// }
	db, err := db.New(cfg.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)

	if err != nil {
		log.Panic(err)
	}

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store: store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}