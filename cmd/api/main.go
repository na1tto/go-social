package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/na1tto/go-social/internal/env"
	repository "github.com/na1tto/go-social/internal/store"
)

func main() {

	cfg := serverConfig{
		addr: env.GetString("ADDR", "8080"),
	}
	
	store := repository.NewStorage(nil)

	app := &application{
		config: cfg,
		store: store,
	}
	
	
	mux := app.mount()
	log.Fatal(app.run(mux))
}
