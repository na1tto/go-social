package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/na1tto/go-social/internal/db"
	"github.com/na1tto/go-social/internal/env"
	repository "github.com/na1tto/go-social/internal/store"
)

func main() {

	cfg := serverConfig{
		addr: env.GetString("ADDR", "8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/gosocial?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
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
	log.Printf("database connection pool stablished")

	store := repository.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
