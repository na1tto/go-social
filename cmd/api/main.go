package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/na1tto/go-social/internal/env"
)

func main() {

	cfg := serverConfig{
		addr: env.GetString("ADDR", "8080"),
	}

	app := &application{
		config: cfg,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
