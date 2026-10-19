package main

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/prabej7/social/internal/db"
	"github.com/prabej7/social/internal/env"
	"github.com/prabej7/social/internal/store"
)

const version = "0.0.1"

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

//	@title			LoferSocial API
//	@description	API documentation for LoferSocial, a social network for lofers.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/api/v1
//
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	if err := env.LoadFile(".envrc"); err != nil {
		log.Fatalf("failed to load environment file: %v", err)
	}

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiUrl: env.GetString("EXTERNAL_API", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15min"),
		},
		env: env.GetString("ENV", "development"),
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config:  cfg,
		store:   store,
		timeout: time.Second * 5,
	}

	log.Printf("server is running at: http://localhost%s", app.config.addr)

	log.Fatal(app.run(app.mount()))
}
