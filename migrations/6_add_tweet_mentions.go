package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table tweet_mentions...")

		_, err := db.Exec(`CREATE TABLE tweet_mentions (
				id					integer NOT NULL PRIMARY KEY,
				account_id			bigint,
				tweet_id			varchar(20)
			);

			ALTER TABLE tweet_mentions
				ADD CONSTRAINT fk_tweet_mentions_tweets
				FOREIGN KEY (tweet_id) REFERENCES tweets(id);

			ALTER TABLE tweet_mentions
				ADD CONSTRAINT fk_tweet_mentions_accounts
				FOREIGN KEY (account_id) REFERENCES accounts(id);
		`)

		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table tweet_mentions...")
		_, err := db.Exec(`DROP TABLE tweet_mentions`)
		return err
	})
}
