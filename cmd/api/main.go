package main

import (
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/na1tto/go-social/internal/db"
	"github.com/na1tto/go-social/internal/env"
	"github.com/na1tto/go-social/internal/mailer"
	repository "github.com/na1tto/go-social/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Go-Social API
//	@description	This is my swagger documentation for the Go-Social Project!
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Luiz (na1tto)
//	@contact.url	https://github.com/na1tto
//	@contact.email	eduardo3245.ss@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {

	cfg := serverConfig{
		addr:        env.GetString("ADDR", "8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/gosocial?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			mailTrap: mailTrapConfig{
				apiKey:          env.GetString("MAILTRAP_API_KEY", ""),
				sandboxUsername: env.GetString("MAILTRAP_USERNAME", ""),
				password:        env.GetString("MAILTRAP_PASSWORD", ""),
			},
			sendgrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
	}

	//logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync() //flushes any buffered log entries

	//db
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
	logger.Info("database connection pool stablished")

	store := repository.NewStorage(db)

	// mailer := mailer.NewSendgrid(cfg.mail.sendgrid.apiKey, cfg.mail.fromEmail)
	mailTrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail, cfg.mail.mailTrap.sandboxUsername, cfg.mail.mailTrap.password)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailTrap,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
