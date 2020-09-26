package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table countries...")

		_, err := db.Exec(`
			CREATE TABLE countries (
				id					integer NOT NULL PRIMARY KEY,
				code				varchar(2),
				name 				varchar(200)	
			);
	
			ALTER TABLE tweets
				ADD CONSTRAINT fk_countries_tweets
				FOREIGN KEY (country_id) REFERENCES countries(id);
		`)

		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table countries...")
		_, err := db.Exec(`
			ALTER TABLE tweets DROP FOREIGN KEY fk_countries_tweets;
			DROP TABLE countries;
		`)
		return err
	})
}
