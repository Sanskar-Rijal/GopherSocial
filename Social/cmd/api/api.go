package main

import (
	"fmt"
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
	db dbConfig
	env string
}

type dbConfig struct {
	addr string
	maxOpenConns int 
	maxIdleConns int 
	maxIdleTime string
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


	//when no route matches then 
	router.NotFound(func(w http.ResponseWriter , r *http.Request){
		app.notFoundError(w,r,fmt.Errorf("Route not found %s", r.URL.Path))
	} )

	//Method Not Allowed - route exists but method is not allowed 
	// example route is GET /posts
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request){
		app.methodNotAllowedError(w,r,fmt.Errorf("method not allowed %s", r.Method))
	})

	router.Route("/v1", func (r chi.Router){
		r.Get("/health",app.healthCheckHandler)

		//creating request for post (add, delete, edit etc)
		r.Route("/posts", func ( r chi.Router){
			//create post
			r.Post("/",app.createPostHandler)
			//get post by id 
			r.Route("/{postID}", func(r chi.Router){
				//using our middleares 
				r.Use(app.postContextMiddleware)
				//routes
				r.Get("/", app.getPostByIdHandler)
				r.Delete("/",app.deletePostHandler)
				r.Patch("/",app.updatePostHandler)
			})

		})
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