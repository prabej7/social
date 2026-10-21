package main

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/prabej7/social/internal/db"
	"github.com/prabej7/social/internal/env"
	"github.com/prabej7/social/internal/mailer"
	"github.com/prabej7/social/internal/store"
	"go.uber.org/zap"
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
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGridConfig: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			smtpConfig: smtpConfig{
				password: env.GetString("SMTP_APP_PASSWORD", ""),
			},
		},
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
	}

	//logger
	logger := zap.Must(zap.NewProduction()).Sugar()

	// database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	mailer := mailer.NewSMTPMailer(cfg.mail.fromEmail, cfg.mail.smtpConfig.password)
	// mailer := mailer.NewSendGrid(cfg.mail.sendGridConfig.apiKey, cfg.mail.fromEmail)
	store := store.NewStorage(db)

	app := &application{
		config:  cfg,
		store:   store,
		timeout: time.Second * 5,
		logger:  logger,
		mailer:  mailer,
	}

	logger.Infof("server is running at: http://localhost%s", app.config.addr)

	logger.Fatal(app.run(app.mount()))
}
