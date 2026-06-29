package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"social/docs"
	"social/internal/auth"
	"social/internal/mailer"
	"social/internal/store"
	"social/internal/store/cache"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
}

type config struct {
	addr   string //address of the port we want to listen on
	db     dbConfig
	env    string
	apiURL string
	mail   mailConfig
	auth   authConfig
	redisCfg redisConfig
}

type redisConfig struct {
	addr string 
	pw string
	db int 
	enabled bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type basicConfig struct {
	username string
	password string
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
	aud    string
}

type mailConfig struct {
	fromEmail     string
	password      string
	isDevelopment bool
	exp           time.Duration
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
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
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFoundError(w, r, fmt.Errorf("Route not found %s", r.URL.Path))
	})

	//Method Not Allowed - route exists but method is not allowed
	// example route is GET /posts
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.methodNotAllowedError(w, r, fmt.Errorf("method not allowed %s", r.Method))
	})

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)

		//Basic authentication
		r.Group(func(r chi.Router) {
			//using middleware
			r.Use(app.BasicAuthMiddleware())
			r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		})

		//Public Routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			//login user
			r.Post("/login", app.LoginUserHandler)
		})

		//comments
		r.Route("/comments", func(r chi.Router) {
			//protect middleware to protect posts 
			r.Use(app.protect)
			r.Post("/", app.addCommentHandler)

			r.Route("/{commentID}", func(r chi.Router) {
				//delete comment
				r.Delete("/", app.deleteCommentHandler)
			})
		})

		//creating request for post (add, delete, edit etc)
		r.Route("/posts", func(r chi.Router) {
			//protect middleware to protect posts 
			r.Use(app.protect)
			//create post
			r.Post("/", app.createPostHandler)
			//get post by id
			r.Route("/{postID}", func(r chi.Router) {
				//using our middleares
				r.Use(app.postContextMiddleware)
				//routes
				r.Get("/", app.getPostByIdHandler)
				r.Delete("/", app.checkPostOwnerShip("admin",app.deletePostHandler))
				r.Patch("/", app.checkPostOwnerShip("moderator" ,app.updatePostHandler))
			})
		})
		//users
		r.Route("/users", func(r chi.Router) {
			//Activate user
			r.Post("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				//using middleware to fetch user from url
				//protect middleware to protect posts 
				r.Use(app.protect)
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserByIdHandler)
				//Route to follow users /v1/users/{userID}/follow
				r.Post("/follow", app.followUserHandler)
				r.Post("/unfollow", app.unfollowUserHandler)

				r.Get("/getfollowers", app.getFollowersHandler)
				r.Get("/getfollowing", app.getFollowingHandler)
			})

			//user feed
			r.Group(func(r chi.Router) {
				r.Use(app.protect)
				r.Get("/feed", app.getUserFeedHandler) //"/v1/users/feed"
			})
		})
	})

	return router
}

func (app *application) run(mux http.Handler) error {

	//Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	//GraceFul server Shutdown 
	// channel that carries the shutdown error
    // main goroutine waits on this
	shutdown := make(chan error)
	//           ↑
    //    create the walkie talkie pipe

	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

    // launch a GOROUTINE — runs in background
    // main goroutine continues below
	go func(){

		// channel that receives OS signals
        // like ctrl+c
		quit := make(chan os.Signal, 1)
		//                          ↑
        //                    buffer of 1 — dont miss the signal

		 // tell Go which signals to listen for
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		 //                    ↑               ↑
        //                 ctrl+c          kill command from terminal
        //                                 docker stop sends this

		// BLOCK here — wait for signal
        // does nothing until ctrl+c or kill is received
		s := <- quit
		//      ↑
        //   waiting... waiting... waiting...
        //   user presses ctrl+c
        //   s = the signal received

	    // give server 5 seconds to finish current requests
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		// tell server to shutdown gracefully
        // waits for current requests to finish (max 5 seconds)
        // then sends result into shutdown channel
		shutdown <-  server.Shutdown(ctx)

	}()
	app.logger.Infow("Server has started at port", "addr", app.config.addr, "env", app.config.env)

	// start server — blocks here
    // server runs until Shutdown() is called
	err := server.ListenAndServe()

	// when Shutdown() is called
    // ListenAndServe returns http.ErrServerClosed
    // that is NORMAL — not a real error
	if !errors.Is(err, http.ErrServerClosed){
		return err
	}
	 // if ErrServerClosed — continue below

	// wait for shutdown goroutine to finish
    // BLOCKS here until goroutine sends result
	err = <-shutdown

	 //      ↑
    //  waiting for goroutine to send result
    //  goroutine sends: shutdown <- server.Shutdown(ctx)
    //  this unblocks

	if err != nil {
		return err
	}

	// return server.ListenAndServe()

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil 
}
