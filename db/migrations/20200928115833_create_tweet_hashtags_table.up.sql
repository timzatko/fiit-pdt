CREATE TABLE tweet_hashtags (
    id					SERIAL PRIMARY KEY,
    hashtag_id			integer,
    tweet_id			varchar(20),
    UNIQUE (hashtag_id, tweet_id)
);

ALTER TABLE tweet_hashtags
    ADD CONSTRAINT fk_tweet_hashtags_tweets
        FOREIGN KEY (tweet_id) REFERENCES tweets(id);

ALTER TABLE tweet_hashtags
    ADD CONSTRAINT fk_tweet_hashtags_hashtags
        FOREIGN KEY (hashtag_id) REFERENCES hashtags(id);