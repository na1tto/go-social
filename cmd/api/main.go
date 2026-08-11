package main

import (
	"expvar"
	"runtime"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/na1tto/go-social/internal/auth"
	"github.com/na1tto/go-social/internal/db"
	"github.com/na1tto/go-social/internal/env"
	"github.com/na1tto/go-social/internal/mailer"
	rateLimiter "github.com/na1tto/go-social/internal/ratelimiter"
	repository "github.com/na1tto/go-social/internal/store"
	"github.com/na1tto/go-social/internal/store/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "1.2.1"

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
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/gosocial?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
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
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 25 * 3, //3 days
				iss:    "GoSocial",
			},
		},
		rateLimiter: rateLimiter.Config{
			RequestPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame:           time.Second * 3,
			Enabled:             env.GetBool("RATE_LIMITER_ENABLED", true),
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

	// cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("redis cache stablished")
	}

	store := repository.NewStorage(db)
	cacheStore := cache.NewRedisStorage(rdb)

	// mailer := mailer.NewSendgrid(cfg.mail.sendgrid.apiKey, cfg.mail.fromEmail)
	mailTrap, _ := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail, cfg.mail.mailTrap.sandboxUsername, cfg.mail.mailTrap.password)

	// instantiate the authenticator
	jwtAuth := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	// rate limiter
	rateLimiter := rateLimiter.NewFixedWindowRateLimiter(
		cfg.rateLimiter.RequestPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStore,
		logger:        logger,
		mailer:        mailTrap,
		authenticator: jwtAuth,
		rateLimiter:   rateLimiter,
	}

	// Metrics colected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
