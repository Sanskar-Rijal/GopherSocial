package main

import (
	"social/internal/db"
	"social/internal/env"
	"social/internal/store"

	"go.uber.org/zap"
)

type ErrorResponseWrapper struct {
	Status  bool   `json:"status" example:"false"`
	Message string `json:"message"`
	Env     string `json:"env" example:"development"`
}

type SuccessResponse[T any] struct {
	Status  bool   `json:"status" example:"true"`
	Env     string `json:"env" example:"development"`
	Message T      `json:"message"`
}

//	@title			GopherSocial API (Social Media)
//	@description	You can use endpoints and make your own frontend 😘
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

const version = "1.0.0"

func main() {

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://sanskar:adminpassword@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("Go_ENV", "dev"),
	}

	//Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//Database
	// app := &application{
	// 	config: config{
	// 		addr: env.GetString("ADDR",":8080"),
	// 		db: dbConfig{
	// 			addr: env.GetString("DB_ADDR","postgres://cutiee:cutiee123@localhost:5432/social?sslmode=disable"),
	// 			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS",30),
	// 			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS",30),
	// 			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME","15m"),
	// 		},
	// 	},
	// 	store: store,
	// }
	db, err := db.New(cfg.db.addr,
		 cfg.db.maxOpenConns, 
		 cfg.db.maxIdleConns,
		  cfg.db.maxIdleTime)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database Connection is live")

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
