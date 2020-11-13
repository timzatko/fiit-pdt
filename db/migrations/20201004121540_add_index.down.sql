ALTER TABLE tweets DROP CONSTRAINT fk_tweets_tweets;
ALTER TABLE tweets DROP CONSTRAINT fk_accounts_tweets;
ALTER TABLE tweets DROP CONSTRAINT fk_countries_tweets;
ALTER TABLE tweet_hashtags DROP CONSTRAINT fk_tweet_hashtags_tweets;
ALTER TABLE tweet_hashtags DROP CONSTRAINT fk_tweet_hashtags_hashtags;
ALTER TABLE tweet_mentions DROP CONSTRAINT fk_tweet_mentions_tweets;
ALTER TABLE tweet_mentions DROP CONSTRAINT fk_tweet_mentions_accounts;
