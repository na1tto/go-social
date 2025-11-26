package main

import (
	"log"

	"github.com/na1tto/go-social/internal/db"
	"github.com/na1tto/go-social/internal/env"
	repository "github.com/na1tto/go-social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/gosocial?sslmode=disable")

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	
	defer conn.Close()
	
	store := repository.NewStorage(conn)
	
	db.Seed(store)
}
