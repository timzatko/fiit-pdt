package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table accounts...")

		_, err := db.Exec(`
			CREATE TABLE accounts (
				id					bigint NOT NULL PRIMARY KEY,
				screen_name			varchar(200),
				name				varchar(200),
				description		 	text,
				followers_count		integer,
				friends_count		integer,
				statuses_count		integer 
			);
	
			ALTER TABLE tweets
				ADD CONSTRAINT fk_accounts_tweets
				FOREIGN KEY (author_id) REFERENCES accounts(id);
		`)

		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table accounts...")
		_, err := db.Exec(`
			ALTER TABLE tweets DROP CONSTRAINT fk_accounts_tweets;
			DROP TABLE accounts;
		`)
		return err
	})
}
