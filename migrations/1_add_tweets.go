package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table tweets...")

		_, err := db.Exec(`CREATE TABLE tweets (
    		id				varchar(20) NOT NULL,
			content			text,
			location 		geometry(point, 4326),
  			retweet_count	integer,
			favorite_count	integer,
			happened_at		timestamptz,
			author_id		bigint,
			country_id		integer,
			parent_id		varchar(20),
			PRIMARY KEY (id),
    		FOREIGN KEY (parent_id) REFERENCES tweets(id)
		)`)

		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table tweets...")
		_, err := db.Exec(`DROP TABLE tweets`)
		return err
	})
}
