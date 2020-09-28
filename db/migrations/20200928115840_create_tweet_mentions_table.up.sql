CREATE TABLE tweet_mentions (
    id					SERIAL PRIMARY KEY,
    account_id			bigint,
    tweet_id			varchar(20),
    UNIQUE (account_id, tweet_id)
);

ALTER TABLE tweet_mentions
    ADD CONSTRAINT fk_tweet_mentions_tweets
        FOREIGN KEY (tweet_id) REFERENCES tweets(id);

ALTER TABLE tweet_mentions
    ADD CONSTRAINT fk_tweet_mentions_accounts
        FOREIGN KEY (account_id) REFERENCES accounts(id);