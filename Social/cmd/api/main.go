package main

import (
	"log"
	"social/internal/env"
	"social/internal/store"
)

func main() {

	store := store.NewPostgresStorage(nil)

	app := &application{
		config: config{
			addr: env.GetString("ADDR",":8080"),
		},
		store: store,
	}


	mux := app.mount()
	log.Fatal(app.run(mux))
}