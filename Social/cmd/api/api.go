package main

import (
	"log"
	"net/http"
	"social/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config 
	store store.Storage
}

type config struct {
	addr string //address of the port we want to listen on
}


func (app *application) mount() http.Handler {
	// mux := http.NewServeMux()

	// mux.HandleFunc("GET /v1/health", app.healthCheckHandler)
	router := chi.NewRouter()


// A good base middleware stack
    router.Use(middleware.RequestID)
    router.Use(middleware.RealIP)
    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)

	router.Route("/v1", func (r chi.Router){
		r.Get("/health",app.healthCheckHandler)
	})


	return router
}

func (app *application) run(mux http.Handler) error {

	server := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second *10,
		IdleTimeout: time.Minute,
	}
	log.Printf("Server has started at %s",app.config.addr)
	return server.ListenAndServe()
}