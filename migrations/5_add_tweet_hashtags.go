package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table tweet_hashtags...")

		_, err := db.Exec(`CREATE TABLE tweet_hashtags (
				id					integer NOT NULL PRIMARY KEY,
				hashtag_id			integer,
				tweet_id			varchar(20)
			);

			ALTER TABLE tweet_hashtags
				ADD CONSTRAINT fk_tweet_hashtags_tweets
				FOREIGN KEY (tweet_id) REFERENCES tweets(id);

			ALTER TABLE tweet_hashtags
				ADD CONSTRAINT fk_tweet_hashtags_hashtags
				FOREIGN KEY (hashtag_id) REFERENCES hashtags(id);
		`)

		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table tweet_hashtags...")
		_, err := db.Exec(`DROP TABLE tweet_hashtags`)
		return err
	})
}
