package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table hashtags...")

		_, err := db.Exec(`CREATE TABLE hashtags (
    		id				integer NOT NULL PRIMARY KEY,
			value			text
		)`)

		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table hashtags...")
		_, err := db.Exec(`DROP TABLE hashtags`)
		return err
	})
}
