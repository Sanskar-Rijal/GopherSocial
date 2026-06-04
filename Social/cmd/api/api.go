package main

import (
	"fmt"
	"net/http"
	"social/docs"
	"social/internal/auth"
	"social/internal/mailer"
	"social/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
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
			r.Post("/", app.addCommentHandler)

			r.Route("/{commentID}", func(r chi.Router) {
				//delete comment
				r.Delete("/", app.deleteCommentHandler)
			})
		})

		//creating request for post (add, delete, edit etc)
		r.Route("/posts", func(r chi.Router) {
			//create post
			r.Post("/", app.createPostHandler)
			//get post by id
			r.Route("/{postID}", func(r chi.Router) {
				//using our middleares
				r.Use(app.postContextMiddleware)
				//routes
				r.Get("/", app.getPostByIdHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})
		//users
		r.Route("/users", func(r chi.Router) {
			//Activate user
			r.Post("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				//using middleware to fetch user from url
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

	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	app.logger.Infow("Server has started at port", "addr", app.config.addr, "env", app.config.env)
	return server.ListenAndServe()
}
