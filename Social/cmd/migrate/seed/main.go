package main

import (
	"log"
	"social/internal/db"
	"social/internal/env"
	"social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://sanskar:adminpassword@localhost:5432/social?sslmode=disable")

	conn, err := db.New(addr, 10, 10, "15m")

	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	store := store.NewPostgresStorage(conn)

	db.Seed(store)
}
