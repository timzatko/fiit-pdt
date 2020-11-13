CREATE TABLE tweet_hashtags (
    id					SERIAL PRIMARY KEY,
    hashtag_id			integer,
    tweet_id			varchar(20),
    UNIQUE (hashtag_id, tweet_id)
);
