CREATE TABLE tweet_mentions (
    id					SERIAL PRIMARY KEY,
    account_id			bigint,
    tweet_id			varchar(20),
    UNIQUE (account_id, tweet_id)
);
