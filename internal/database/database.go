package database

import (
	"log"

	"github.com/go-pg/pg"
)

func Connect() *pg.DB {
	return pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Password: "123",
		Database: "tweets",
	})
}

func Close(db *pg.DB) {
	err := db.Close()

	if err != nil {
		log.Print("error while closing the database...")
		log.Fatal(err)
	}
}
